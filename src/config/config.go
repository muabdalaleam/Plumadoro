package config

import "time"

type Config struct {
	tickDuration   time.Duration
	padding        uint16
	maxWidth       uint16
}


