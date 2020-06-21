package config

import "time"

type Config struct {
	PomodoroDuration    *time.Duration
	ShortBreaksDuration *time.Duration
	LongBreaksDuration  *time.Duration
}

func NewConfig() *Config {
	return &Config{}
}
