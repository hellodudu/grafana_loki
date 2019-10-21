package loki_conn

import (
	"time"
)

type Options struct {
	URL      string        `flag:"url"`
	Interval time.Duration `flag:"interval"`
}

func NewOptions() *Options {

	return &Options{
		URL:      "loki:3100/api/prom/push",
		Interval: 1 * time.Minute,
	}
}
