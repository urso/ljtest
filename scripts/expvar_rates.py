import requests
import argparse
import time
import curses


def main():
    parser = argparse.ArgumentParser(
        description="Print per second stats from expvars")
    parser.add_argument("url",
                        help="The URL from where to read the values")
    args = parser.parse_args()

    stdscr = curses.initscr()
    curses.noecho()
    curses.cbreak()

    last_vals = {}

    # running average for last 30 measurements
    N = 30
    avg_vals = {}
    now = time.time()

    while True:
        try:
            time.sleep(1.0)
            stdscr.erase()

            r = requests.get(args.url)
            json = r.json()

            last = now
            now = time.time()
            dt = now - last

            for key, total in json.items():
                if isinstance(total, (int, long, float)):
                    if key in last_vals:
                        per_sec = (total - last_vals[key])/dt
                        if key not in avg_vals:
                            avg_vals[key] = []
                        avg_vals[key].append(per_sec)
                        if len(avg_vals[key]) > N:
                            avg_vals[key] = avg_vals[key][1:]
                        avg_sec = sum(avg_vals[key])/len(avg_vals[key])
                    else:
                        per_sec = "na"
                        avg_sec = "na"
                    last_vals[key] = total
                    try:
                        stdscr.addstr("{}: {}/s (avg: {}/s) (total: {})\n"
                                      .format(key, per_sec, avg_sec, total))
                    except:
                        pass
            stdscr.refresh()
        except requests.ConnectionError:
            stdscr.addstr("Waiting for connection...\n")
            stdscr.refresh()
            last_vals = {}
            avg_vals = {}

if __name__ == "__main__":
    main()
