package main

import "time"

type Options struct {
	Duration time.Duration `config:"opt.duration"`
}

var (
	defaultOptions = Options{
		Duration: 20 * time.Second,
	}
)
