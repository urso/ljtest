package main

import "time"

type Options struct {
	Duration  time.Duration `config:"duration"`
	WaitAfter time.Duration `config:"wait.after"`
}

var (
	defaultOptions = Options{
		Duration: 20 * time.Second,
	}
)
