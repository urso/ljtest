package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/urso/ucfg"
	"github.com/urso/ucfg/cfgutil"
	cfgflag "github.com/urso/ucfg/flag"
	"github.com/urso/ucfg/yaml"
)

func main() {
	dir := flag.String("d", "config", "Runner configurations")
	tmpl := flag.String("t", "templates", "Template directory")
	cliConfig := cfgflag.Config(nil, "D", "Configuration overwrites",
		ucfg.VarExp, ucfg.PathSep("."))
	flag.Parse()

	cfg, err := readConfig(*dir, *tmpl, cliConfig, flag.Args())
	if err != nil {
		log.Println("Failed to load configuration")
		log.Println(err)
		os.Exit(1)
	}

	if err := run(cfg); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func run(cfg *ucfg.Config) error {
	runs := struct {
		Runs map[string]*ucfg.Config `config:"run"`
		Opt  Options                 `config:"opt"`
	}{nil, defaultOptions}

	err := cfg.Unpack(&runs, ucfg.PathSep("."), ucfg.VarExp, ucfg.ResolveEnv)
	if err != nil {
		log.Println("Failed to unpack run configuration")
		return err
	}

	for name, active := range runs.Runs {
		log.Println("start: ", name)

		r := runConfig{}
		if err := active.Unpack(&r); err != nil {
			return err
		}

		_active, _ := ucfg.NewFrom(map[string]interface{}{
			"name": name,
			"run":  active,
		}, ucfg.PathSep("."), ucfg.VarExp)
		scenarios, err := makeScenarios(cfg, r, _active)
		if err != nil {
			return err
		}

		for i, s := range scenarios {
			log.Println("  run: ", i)
			if err := s.run(&runs.Opt, cfg); err != nil {
				log.Println("  failed")
				return err
			}

			if runs.Opt.WaitAfter > 0 {
				log.Println("waiting: ", runs.Opt.WaitAfter)
				time.Sleep(runs.Opt.WaitAfter)
			}
		}
		log.Println("finished: ", name)
	}

	return nil
}

func readConfig(
	configDir string,
	tmplDir string,
	cliConfig *ucfg.Config,
	runFiles []string,
) (*ucfg.Config, error) {
	workdir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cfgOpts := []ucfg.Option{ucfg.PathSep("."), ucfg.VarExp}

	cfg := cfgutil.NewCollector(nil, cfgOpts...)
	cfg.Add(ucfg.NewFrom(map[string]interface{}{
		"path.config":   configDir,
		"path.template": tmplDir,
		"path.workdir":  workdir,
	}, cfgOpts...))

	// load configs
	err = filepath.Walk(configDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Ext(path) != ".yml" || filepath.Base(path)[0] == '.' {
			return err
		}
		return cfg.Add(yaml.NewConfigWithFile(path, cfgOpts...))
	})
	if err != nil {
		return nil, err
	}

	for _, path := range runFiles {
		if err := cfg.Add(yaml.NewConfigWithFile(path, cfgOpts...)); err != nil {
			return nil, err
		}
	}

	cfg.Add(cliConfig, nil)
	return cfg.Get()
}
