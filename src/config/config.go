package config

import (
	"fmt"
	"os"
	"time"
	"errors"
	toml "github.com/BurntSushi/toml"
)

type Config struct {
	TickDuration   time.Duration `toml:"tick_duration"`
	Padding        uint16        `toml:"padding"`
	MaxWidth       uint16        `toml:"max_width"`
}

// TODO: add diffrenet config paths to search throw.
const configPath = "${HOME}/.plumadoro.toml"

var defaultConfig = Config {
	TickDuration   : time.Millisecond * 20,
	Padding        : 10,
	MaxWidth       : 40,
}

var (
	ErrUnsupportedKeys     error = errors.New("Unsupported keys in the TOML config")
	ErrFailedReadingConfig error = errors.New("Failed Reading the configuration path")
	ErrFailedParsingTOML   error = errors.New("Failed parsing the TOML configuration file")
)

func LoadConfig() (Config, error) {
	var err     error  = nil
	var config  Config = defaultConfig

	dat, err := os.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("%w: %w", ErrFailedReadingConfig, err)
	}

	md, err := toml.Decode(string(dat), &config)
	if err != nil {
		return config, fmt.Errorf("%w: %w", ErrFailedParsingTOML, err)
	}

	if undecoded := md.Undecoded(); len(undecoded) != 0 {
		err = fmt.Errorf("%w: Unsupported keys: %q", ErrUnsupportedKeys, undecoded)
	}

	if md.IsDefined("tick_duration") {
		if config.TickDuration < time.Millisecond {
			config.TickDuration *= time.Millisecond
		}  // HACK: ?
	}

	if md.IsDefined("max_width") {
		if config.MaxWidth > 150 {
			err = fmt.Errorf("%w\nInvalid max_width value it must be between 0 and 150.", err)
		}
	}

	if md.IsDefined("padding") {
		if config.Padding > 50 {
			err = fmt.Errorf("%w\nInvalid padding value it must be between 0 and 50.", err)
		}
	}

	return config, err
}

