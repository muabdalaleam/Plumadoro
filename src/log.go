package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"errors"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/progress"
)

type LogTickMsg time.Time

type pomodoroRecord struct {
	remainingTime    time.Duration
	pausedTime       time.Duration
	phaseType        phaseType
	n                uint64
	running          bool
	time_            time.Time
}

var (
	ErrStateNotRestorable  = errors.New("Cannot restore last state because it's not in the same day")
	ErrFailedReadingLog    = errors.New("Failed reading the configuration path no log file found")
	ErrFailedParsingLog    = errors.New("Failed parsing a CSV row in the log file")
)

var cacheDir, _ = os.UserCacheDir()

var logPath string = fmt.Sprintf("%s/.plumadoro_log.csv", cacheDir)

const timeFormat string = time.RFC3339

const LogTickDuration time.Duration = time.Second * 30 // XXX: hardcoded i know


func fromCSVRow(row []string) (pomodoroRecord, error) {
	var record pomodoroRecord
	var err error

	// 6 is the count of pomodoroRecord's fields
	if len(row) != 6 {
		return record, fmt.Errorf("%w: Invalid length for row it must be 6 cols only.", ErrFailedParsingLog)
	}

	var phase phaseType
	switch (row[2]) {
	case "focus"      :  phase = Focus
	case "short_break":  phase = ShortBreak
	case "long_break" :  phase = LongBreak
	}

	record.remainingTime, err = time.ParseDuration(row[0])
	record.pausedTime, err    = time.ParseDuration(row[1])
	record.phaseType          = phase            // row[2]
	record.n, err             = strconv.ParseUint(row[3], 10, 8)
	record.running, err       = strconv.ParseBool(row[4])
	record.time_, err         = time.Parse(timeFormat, row[5])

	if err != nil {
		return record, ErrFailedParsingLog
	}

	return record, nil
}

func (r pomodoroRecord) toCSVRow() []string {
	var phaseTypeStr string
	switch (r.phaseType) {
	case Focus:         phaseTypeStr = "focus"
	case ShortBreak:    phaseTypeStr = "short_break"
	case LongBreak:     phaseTypeStr = "long_break"
	}

	return []string{
		r.remainingTime.Round(time.Second).String(), // Remaning time
		r.pausedTime.Round(time.Second).String(),    // Paused time
		phaseTypeStr,                                // Phase type
		strconv.FormatUint(r.n, 10),                                 // N of the current phase
		strconv.FormatBool(r.running),               // Running
		r.time_.Format(timeFormat),
	}
}

// Should only be used twice when the pomodor app terminates and when a phase ends
func (p *PomodoroModel) save() error {
	var err error = nil
	var file *os.File

	// HACK: IDK how the hell this works but i don't give a s***
	file, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return ErrFailedReadingLog
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write(pomodoroRecord{
		remainingTime: p.remainingTime,
		pausedTime:    p.pausedTime,
		phaseType:     p.phaseType,
		n:             uint64(p.n),
		running:       p.running,
		time_:         time.Now(),
		}.toCSVRow(),
	)

	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedParsingLog, err)
	}

	return err
}

// TODO: use Seek and stat to read the last line only instead of the whole file
func (p *PomodoroModel) restore() error {
	file, err := os.Open(logPath)
	if err != nil {
		return ErrFailedReadingLog
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedParsingLog, err)
	}

	if len(records) == 0 {
		return fmt.Errorf("%w: %w", ErrFailedParsingLog, "Log is empty.")
	}

	record, err := fromCSVRow(records[len(records) - 1])
	if err != nil {
		return err
	}

	now := time.Now()
	year, month, day := record.time_.Date()

	if year != now.Year() || month != now.Month() || day != now.Day() {
		return ErrStateNotRestorable
	}

	p.remainingTime     = record.remainingTime
	p.pausedTime        = record.pausedTime
	p.phaseType         = record.phaseType
	p.running           = Config.Autostart
	p.n                 = uint8(record.n)
	p.phasesDurations   = map[phaseType]time.Duration{
		Focus:      Config.Durations.Focus,
		ShortBreak: Config.Durations.ShortBreak,
		LongBreak:  Config.Durations.LongBreak,
	}
	p.progressBar       = progress.New(
		progress.WithSolidFill(p.getPhaseColor()),
	)

	return nil
}

