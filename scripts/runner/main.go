package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/urso/ucfg"
	"github.com/urso/ucfg/yaml"
)

func main() {
	dir := flag.String("d", "config", "Runner configurations")
	tmpl := flag.String("t", "templates", "Template directory")
	flag.Parse()

	cfg, err := readConfig(*dir, *tmpl, flag.Args())
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
		Opt  Options
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
		}
		log.Println("finished: ", name)
	}

	return nil
}

func readConfig(configDir string, tmplDir string, runFiles []string) (*ucfg.Config, error) {

	workdir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cfg, err := ucfg.NewFrom(map[string]interface{}{
		"path.config":   configDir,
		"path.template": tmplDir,
		"path.workdir":  workdir,
	}, ucfg.PathSep("."), ucfg.VarExp)

	// load configs
	err = filepath.Walk(configDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Ext(path) != ".yml" || filepath.Base(path)[0] == '.' {
			return err
		}
		return mergeConfig(cfg, path)
	})
	if err != nil {
		return nil, err
	}

	for _, runFile := range runFiles {
		if err := mergeConfig(cfg, runFile); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func mergeConfig(to *ucfg.Config, path string) error {
	from, err := yaml.NewConfigWithFile(path, ucfg.PathSep("."), ucfg.VarExp)
	if err != nil {
		return err
	}

	err = to.Merge(from)
	/*
		{
			// print config
			var tmp map[string]interface{}
			// cfg.Unpack(&tmp, ucfg.PathSep("."))
			to.Unpack(&tmp)
			raw, err := json.MarshalIndent(tmp, "> ", "    ")
			log.Println(tmp)
			log.Println(string(raw))
			log.Println(err)
		}
	*/

	return err
}
