package main

import (
	"time"

	"github.com/urso/ucfg"
)

type runConfig struct {
	Services       []serviceRawConfig `config:"services" validate:"required"`
	Configurations []*ucfg.Config     `config:"configurations"`
}

type scenario struct {
	Services []serviceRawConfig
	Config   *ucfg.Config
	Active   *ucfg.Config
}

func makeScenarios(env *ucfg.Config, run runConfig, activeRun *ucfg.Config) ([]scenario, error) {
	var scenarios []scenario
	for _, config := range run.Configurations {
		active, _ := ucfg.NewFrom(activeRun, ucfg.PathSep("."), ucfg.VarExp)
		active.SetChild("config", -1, config)
		active, _ = ucfg.NewFrom(map[string]interface{}{"_active": active})

		check := struct {
			Enabled bool `config:"enable"`
		}{true}
		err := config.Unpack(&check,
			ucfg.PathSep("."), ucfg.Env(env), ucfg.Env(active), ucfg.ResolveEnv)
		if err != nil {
			return nil, err
		}

		if !check.Enabled {
			continue
		}

		scenarios = append(scenarios, scenario{
			Services: run.Services,
			Config:   config,
			Active:   active,
		})
	}

	if len(scenarios) == 0 {
		active, _ := ucfg.NewFrom(map[string]interface{}{"_active": activeRun})
		scenarios = []scenario{{Services: run.Services, Active: active}}
	}

	return scenarios, nil
}

func (s *scenario) run(opts *Options, env *ucfg.Config) error {
	services, err := s.readServiceConfig(env)
	if err != nil || len(services) == 0 {
		return err
	}

	for _, service := range services {
		if err := service.run(); err != nil {
			return err
		}
		defer service.stop()
	}
	time.Sleep(opts.Duration)

	return nil
}

func (s *scenario) readServiceConfig(env *ucfg.Config) ([]*service, error) {
	var services []*service

	opts := []ucfg.Option{ucfg.Env(env), ucfg.PathSep("."), ucfg.VarExp, ucfg.ResolveEnv}
	if s.Active != nil {
		opts = append(opts, ucfg.Env(s.Active))
	}

	for _, raw := range s.Services {
		cfg, err := s.mergeServiceConfig(raw, env)
		if err != nil {
			return nil, err
		}

		// unpack shared service configs
		scfg := defaultServiceConfig
		if err := cfg.Unpack(&scfg, opts...); err != nil {
			return nil, err
		}

		if !scfg.Enabled {
			continue
		}

		var raw map[string]interface{}
		cfg.Unpack(&raw, opts...)
		services = append(services, &service{config: &scfg, all: raw})
	}

	return services, nil
}

func (s *scenario) mergeServiceConfig(
	raw serviceRawConfig,
	env *ucfg.Config,
) (*ucfg.Config, error) {
	cfg, err := ucfg.NewFrom(map[string]interface{}{
		"config": map[string]interface{}{},
	}, ucfg.PathSep("."))

	opts := []ucfg.Option{ucfg.Env(env), ucfg.PathSep("."), ucfg.VarExp}

	if err := cfg.Merge(raw.Service, opts...); err != nil {
		return nil, err
	}

	if raw.Config != nil {
		sub, _ := cfg.Child("config", -1)
		if err := sub.Merge(raw.Config, opts...); err != nil {
			return nil, err
		}
	}

	if s.Config != nil {
		sub, _ := cfg.Child("config", -1)
		err = sub.Merge(s.Config, opts...)
		if err != nil {
			return nil, err
		}
	}

	return cfg, err
}
