package main

import (
	"fmt"
	"os"
	"time"
	"errors"
	"strconv"
	"golang.org/x/exp/constraints"
	toml "github.com/BurntSushi/toml"
)

type (
	ConfigT struct {
		TickDuration       time.Duration   `toml:"tick_duration"`
		MaxPauseDuration   time.Duration   `toml:"max_pause_duration"` // NOTE: this max is per phase
		Autostart          bool            `toml:"auto_start"`
		Skipping           bool            `toml:"skipping"` // Allow skipping for phases

		ProgressBar         ProgressBarConfigT  `toml:"progress_bar"`
		Durations           DurationsConfigT    `toml:"durations"`

		loadedConfig        bool // was LoadConfig called before
		loadedConfigPath    string
	}

	ProgressBarConfigT struct  {
		// TODO: Move style options out from here.
		// TODO: support gradient progress bar
		Padding          uint16        `toml:"padding"`
		MaxWidth         uint16        `toml:"max_width"`
		Border           string        `toml:"border_type"`
		FocusColor       string        `toml:"focus_color"`
		ShortBreakColor  string        `toml:"short_break_color"` 
		LongBreakColor   string        `toml:"long_break_color"` 
		PauseColor       string        `toml:"pause_color"` 
		FocusMsg         string        `toml:"focus_msg"`
		ShortBreakMsg    string        `toml:"short_break_msg"`
		LongBreakMsg     string        `toml:"long_break_msg"`
		PauseMsg         string        `toml:"pause_msg"`
	}

	DurationsConfigT struct {
		Focus        time.Duration   `toml:"focus"`
		ShortBreak   time.Duration   `toml:"short_break"`
		LongBreak    time.Duration   `toml:"long_break"`
	}
)

// TODO: handle the case of not finding the home dir or the config dir
var homeDir, _   = os.UserHomeDir()
var configDir, _ = os.UserConfigDir()

var configPaths = []string{
	fmt.Sprintf("%s/.plumadoro.toml", homeDir),
	fmt.Sprintf("%s/.pluma.toml", homeDir),
	fmt.Sprintf("%s/plumadoro.conf", configDir),
	fmt.Sprintf("%s/pluma.conf", configDir),
}
 
var defaultConfig = ConfigT {
	TickDuration      :  time.Millisecond * 20,
	MaxPauseDuration  :  time.Minute * 5,
	Autostart         :  true,
	Skipping          :  false,

	ProgressBar: ProgressBarConfigT{
		Padding   : 5,
		MaxWidth  : 20,

		Border          : "thick",
		FocusColor      : "red",
		ShortBreakColor : "green",
		LongBreakColor  : "cyan",
		PauseColor      : "black",

		FocusMsg        : "Let's Focus",
		ShortBreakMsg   : "Don't watch brainrot",
		LongBreakMsg    : "You deserve it buddy",
		PauseMsg        : "Stop being lazy, jerk",
	},

	Durations: DurationsConfigT{
		Focus       : 25 * time.Minute,
		ShortBreak  : 5 * time.Minute,
		LongBreak   : 20 * time.Minute,
	},

	loadedConfigPath : "",
	loadedConfig:      false,
}

var (
	ErrConfigAlreadyLoaded = errors.New("Calling LoadConfig() while Config global variable is loaded before")
	ErrFailedReadingConfig = errors.New("Failed Reading the configuration paths no config file found")
	ErrFailedParsingTOML   = errors.New("Failed parsing the TOML configuration file")
	ErrUnsupportedKeys     = errors.New("Unsupported keys in the TOML config")
	ErrInvalidKeyValue     = errors.New("Invalid value for TOML key/s") // it parsed well but the value is wrong
)

var Config ConfigT = defaultConfig


func validateRange[T constraints.Ordered](errsPtr *[]error, valuePtr *T, min T, max T, defaultValue T, key string) {
	if *valuePtr < min || *valuePtr > max {
		*valuePtr = defaultValue
		*errsPtr  = append(*errsPtr, fmt.Errorf("Invalid %s: must be between %v and %v", key, min, max))
	}
}

func validateStringLen(errsPtr *[]error, valuePtr *string, min int, max int, defaultValue string, key string) {
	if len(*valuePtr) > max || len(*valuePtr) < min {
		*errsPtr = append(*errsPtr, fmt.Errorf("Invalide %s, it must be between %d, and %d.", key, min, max))
		*valuePtr = defaultValue
	}
}

func validateOption(errsPtr *[]error, valuePtr *string, optionsPtr *[]string, defaultValue string, key string) {
	for _, d := range *optionsPtr {
		if d == *valuePtr {
			return
		}
	}

	*errsPtr = append(*errsPtr, fmt.Errorf("Invalid %s: must be from the following:\n%#v", key, *optionsPtr))
}

func validateColor(errsPtr *[]error, valuePtr *string, defaultValue string, key string) {
	var err error
	value := *valuePtr

	// Based on termenv's docs
	switch (value) {
	case "black":          value = "0"
	case "red":            value = "1"
	case "green":          value = "2"
	case "yellow":         value = "3"
	case "blue":           value = "4"
	case "magenta":        value = "5"
	case "cyan":           value = "6"
	case "white":          value = "7"
	case "bright_black":   value = "8"
	case "bright_red":     value = "9"
	case "bright_green":   value = "10"
	case "bright_yellow":  value = "11"
	case "bright_blue":    value = "12"
	case "bright_magenta": value = "13"
	case "bright_cyan":    value = "14"
	case "bright_white":   value = "15"
	}

	// Checking if value is an ANSI color
	_, err = strconv.ParseUint(value, 10, 8)	
	if err == nil {
		*valuePtr = value
		return
	}

	// Checking if it's a HEX color
	// starting from the second chat because the first one is '#'
	_, err = strconv.ParseUint(value[1:], 16, 32)
	if err == nil {
		return
	}

	// TODO: make this error message more clear
	*valuePtr = defaultValue
	*errsPtr = append(*errsPtr, fmt.Errorf("Invalid %s: must be an ANSI color or HEX color as a string", key))
}

func LoadConfig() error {
	var errs []error // error messages to be concated

	if Config.loadedConfig {
		return ErrConfigAlreadyLoaded
	}
	
	// Loading the config file data
	var dat []byte
	{
		var err error
		for _, configPath := range configPaths {
			dat, err = os.ReadFile(configPath)

			if err == nil {
				Config.loadedConfigPath = configPath
				break
			} // else it will try other paths
		}
		if err != nil {
			return errors.Join(ErrFailedReadingConfig, err)
		}
	}

	// Parsing the TOML config
	{
		md, err := toml.Decode(string(dat), &Config)
		if err != nil {
			return errors.Join(ErrFailedParsingTOML, err)
		}

		if undecoded := md.Undecoded(); len(undecoded) != 0 {
			errs = append(errs, fmt.Errorf("%w: Unsupported keys: %q", ErrUnsupportedKeys, undecoded))
		}
	}

	validateRange(&errs, &Config.TickDuration,
		time.Microsecond, time.Second * 5,
		defaultConfig.TickDuration, "tick_duration")

	// Autostart doesn't require validation

	// Skipping doesn't require validation too.

	validateRange(&errs, &Config.MaxPauseDuration,
		time.Minute*0, time.Minute*1000,
		defaultConfig.MaxPauseDuration, "max_pause_duration")

	validateRange(&errs, &Config.ProgressBar.MaxWidth,
		10, 150,
		defaultConfig.ProgressBar.MaxWidth, "progress_bar.max_width")

	validateRange(&errs, &Config.ProgressBar.Padding,
		0, 50,
		defaultConfig.ProgressBar.Padding, "progress_bar.padding")

	validateOption(&errs, &Config.ProgressBar.Border,
		&[]string{"rounded", "ascii", "thick", "double", "normal", "hidden"}, 
		defaultConfig.ProgressBar.Border, "progress_bar.border_type")

	validateColor(&errs, &Config.ProgressBar.FocusColor,
		defaultConfig.ProgressBar.FocusColor, "progress_bar.focus_color")

	validateColor(&errs, &Config.ProgressBar.ShortBreakColor,
		defaultConfig.ProgressBar.ShortBreakColor, "progress_bar.short_break_color")

	validateColor(&errs, &Config.ProgressBar.LongBreakColor,
		defaultConfig.ProgressBar.LongBreakColor, "progress_bar.long_break_color")

	validateColor(&errs, &Config.ProgressBar.PauseColor,
		defaultConfig.ProgressBar.LongBreakColor, "progress_bar.pause_color")

	validateStringLen(&errs, &Config.ProgressBar.FocusMsg,
		0, 32,
		defaultConfig.ProgressBar.FocusMsg, "progress_bar.focus_msg")

	validateStringLen(&errs, &Config.ProgressBar.ShortBreakMsg,
		0, 32,
		defaultConfig.ProgressBar.ShortBreakMsg, "progress_bar.short_break_msg")

	validateStringLen(&errs, &Config.ProgressBar.LongBreakMsg,
		0, 32,
		defaultConfig.ProgressBar.LongBreakMsg, "progress_bar.long_break_msg")

	validateStringLen(&errs, &Config.ProgressBar.PauseMsg,
		0, 32,
		defaultConfig.ProgressBar.LongBreakMsg, "progress_bar.pause_msg")
	
	validateRange(&errs, &Config.Durations.Focus,
		time.Second*1, time.Minute*1000,
		defaultConfig.Durations.Focus, "duration.focus")

	validateRange(&errs, &Config.Durations.ShortBreak,
		time.Second*1, time.Minute*1000,
		defaultConfig.Durations.ShortBreak, "duration.short_break")

	validateRange(&errs, &Config.Durations.LongBreak,
		time.Second*1, time.Minute*1000,
		defaultConfig.Durations.LongBreak, "duration.long_break")

	err := errors.Join(errs...)

	if err != nil {
		err = fmt.Errorf("%w\n%w", err, ErrInvalidKeyValue)
	}

	return err
}

