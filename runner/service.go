package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/urso/ucfg"
)

type service struct {
	config *serviceConfig
	all    map[string]interface{}

	cmd  *exec.Cmd
	tmpl string
	name string
}

type serviceConfig struct {
	Command        string        `config:"cmd"`
	Name           string        `config:"name"`
	Dir            string        `config:"dir"`
	Enabled        bool          `config:"enabled"`
	Stdout         bool          `config:"stdout"`
	Stderr         bool          `config:"stderr"`
	ConfigFileFlag string        `config:"config.file.flag"`
	ConfigTemplate string        `config:"config.template"`
	WaitBefore     time.Duration `config:"config.wait.before"`
}

type serviceRawConfig struct {
	Service *ucfg.Config `config:"service"`
	Config  *ucfg.Config `config:"config"`
}

var (
	defaultServiceConfig = serviceConfig{
		Enabled: true,
	}
)

func (s *service) run() error {
	cmd, tmpl, err := makeCommand(s.config, s.all)
	if err != nil {
		return err
	}

	s.tmpl = tmpl
	s.cmd = cmd
	s.name = filepath.Base(cmd.Args[0])
	time.Sleep(s.config.WaitBefore)

	log.Println("  start service: ", s.name)
	return s.cmd.Start()
}

func (s *service) stop() {
	if s.cmd != nil {
		s.cmd.Process.Kill()
		s.cmd.Process.Wait()

		s.cmd = nil
		os.Remove(s.tmpl)
	}
}

func makeCommand(service *serviceConfig, config map[string]interface{}) (*exec.Cmd, string, error) {
	args, err := parseCommand(service.Command)
	if err != nil {
		return nil, "", err
	}

	cmd := exec.Command(args[0], args[1:]...)
	if service.Name != "" {
		cmd.Args[0] = service.Name
	}
	cmd.Dir = service.Dir

	tmpl := ""
	if service.ConfigTemplate != "" {
		tmpl, err = loadTemplate(service.ConfigTemplate, config)
		if err != nil {
			return nil, "", err
		}
		log.Printf("  -- new template file (service=%v): %v",
			filepath.Base(cmd.Args[0]), tmpl)

		if service.ConfigFileFlag != "" {
			cmd.Args = append(cmd.Args, service.ConfigFileFlag)
		}
		cmd.Args = append(cmd.Args, tmpl)
	}

	if service.Stdout {
		cmd.Stdout = os.Stdout
	}
	if service.Stderr {
		cmd.Stderr = os.Stderr
	}

	return cmd, tmpl, nil
}

func loadTemplate(path string, cfg map[string]interface{}) (string, error) {
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return "", err
	}

	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}
	defer tmp.Close()

	err = tmpl.Execute(tmp, cfg)
	if err != nil {
		os.Remove(tmp.Name())
	}

	return tmp.Name(), nil
}

func parseCommand(in string) ([]string, error) {
	trim := func(s string) string { return strings.TrimFunc(s, unicode.IsSpace) }

	raw := trim(in)
	if len(raw) == 0 {
		return nil, fmt.Errorf("Empty command given")
	}

	var args []string
	for len(raw) > 0 {
		idx := strings.IndexFunc(raw, func(r rune) bool {
			return r == '"' || unicode.IsSpace(r)
		})
		if idx < 0 {
			args = append(args, trim(raw))
			break
		}

		switch raw[idx] {
		case '"':
			if str := trim(raw[:idx]); str != "" {
				args = append(args, str)
			}
			raw = raw[idx+1:]
			idx = strings.IndexRune(raw, '"')
			if idx < 0 {
				return nil, fmt.Errorf("Failed parsing '%v' due to missing '\"'", in)
			}
			args = append(args, raw[:idx])
			raw = trim(raw[idx+1:])
		default:
			args = append(args, trim(raw[:idx]))
			raw = trim(raw[idx+1:])
		}
	}
	return args, nil
}
