package pomodoro

import (
	"fmt"
	"time"
	"plumadoro/config"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/progress"
)


type TickMsg time.Time

type phaseType byte

// TODO: make a proper error handling for the pomodoro model
type PomodoroModel struct {
	remainingTime    time.Duration
	pausedTime       time.Duration
	phaseType        phaseType
	running          bool
	n                uint8 

	// Configurable
	phasesDurations  map[phaseType]time.Duration
	ProgressBar      progress.Model // HACK: WTH is this
}

const (
	Focus phaseType = iota
	ShortBreak
	LongBreak
)


func (p *PomodoroModel) loadBorder() lipgloss.Style {
	var borderStyle lipgloss.Style

	{
		var borderColor string
		var borderType  lipgloss.Border

		if p.running {
			switch (p.phaseType) {
				case Focus:       borderColor = config.Config.ProgressBar.FocusColor
				case ShortBreak:  borderColor = config.Config.ProgressBar.ShortBreakColor
				case LongBreak:   borderColor = config.Config.ProgressBar.LongBreakColor
			}
		} else {
			borderColor = config.Config.ProgressBar.PauseColor
		}

		switch (config.Config.ProgressBar.Border) {
		case "rounded":   borderType = lipgloss.RoundedBorder()
		case "ascii":     borderType = lipgloss.ASCIIBorder()
		case "thick":     borderType = lipgloss.ThickBorder()
		case "double":    borderType = lipgloss.DoubleBorder()
		case "normal":    borderType = lipgloss.NormalBorder()
		case "hidden":    borderType = lipgloss.HiddenBorder()
		}

		borderStyle = lipgloss.NewStyle().
			Padding(1, 2).
			BorderStyle(borderType).
			BorderForeground(lipgloss.Color(borderColor))
	}

	return borderStyle
}

func (p *PomodoroModel) loadPhaseMsg() string {
	var msg string

	// TODO: add some styling to it
	if p.running {
		switch (p.phaseType) {
		case Focus:        msg = config.Config.ProgressBar.FocusMsg
		case ShortBreak:   msg = config.Config.ProgressBar.ShortBreakMsg
		case LongBreak:    msg = config.Config.ProgressBar.LongBreakMsg
		}
	} else {
		msg = config.Config.ProgressBar.PauseMsg
	}

	return msg
}

func (p *PomodoroModel) loadProgressBarColor() string { 
	var progressColor string
	switch (p.phaseType) {
		case Focus:       progressColor = config.Config.ProgressBar.FocusColor
		case ShortBreak:  progressColor = config.Config.ProgressBar.ShortBreakColor
		case LongBreak:   progressColor = config.Config.ProgressBar.LongBreakColor
	}

	return progressColor
}


// TODO: load the phase from the history
func (p *PomodoroModel) Init() {
	if err := p.restore(); err != nil {
		// XXX: Handle errors and show them to the user somehow.

		*p = PomodoroModel {
			remainingTime: config.Config.Durations.Focus,
			pausedTime:    time.Duration(0),
			phaseType:     Focus,
			running:       config.Config.Autostart,
			n:             1, // NOTE: the index of phases is one based

			phasesDurations: map[phaseType]time.Duration{
				Focus:      config.Config.Durations.Focus,
				ShortBreak: config.Config.Durations.ShortBreak,
				LongBreak:  config.Config.Durations.LongBreak,
			},
			ProgressBar:     progress.New(
				progress.WithSolidFill(config.Config.ProgressBar.FocusColor), // HACK: I hate my life
			),
		}
	}
}

// TODO: add a maximum pause time
func (p *PomodoroModel) Toggle() {
	if p.running == false {
		p.running = true
	} else {
		p.running = false
	}
}

func (p *PomodoroModel) Tick(d time.Duration) {
	if p.running {
		p.remainingTime -= d
	} else {
		p.pausedTime += d
	}

	if p.remainingTime <= time.Duration(0) {
		p.Next()
	}
}

// TODO: Make a reset function
// func (p *PomodoroModel) Reset() {
// }

// It updates the whole state of the PomodoroModel
func (p *PomodoroModel) Next() {
	// NOTE: be careful n is updated first
	p.n += 1

	var newPhaseType phaseType

	switch p.phaseType {
	case ShortBreak, LongBreak:
		newPhaseType = Focus
	// TODO: Checking next phase should be more configurable and not hard coded
	case Focus:
		if p.n == 7 || (p.n - 1) % 7 == 0 { // because for every 4 focus phases there are three short breaks in between
			newPhaseType = LongBreak
		} else {
			newPhaseType = ShortBreak
		}
	}

	p.pausedTime    = time.Duration(0)
	p.remainingTime = p.phasesDurations[newPhaseType]
	p.phaseType     = newPhaseType
	p.running       = config.Config.Autostart
	p.ProgressBar.FullColor = p.loadProgressBarColor() // TODO: add the abilitiy to use gradients

	err := p.save()
	if err != nil {
		panic(err) // XXX: PANIC FOR REAL!
	}
}

func (p *PomodoroModel) Render() string {
	s := lipgloss.JoinVertical(
		lipgloss.Center,
		p.loadBorder().Render(p.ProgressBar.ViewAs(p.GetProgress())),
		"",
		fmt.Sprintf("Remaining: %02d:%02d",
			int(p.remainingTime.Minutes()),
			int(p.remainingTime.Seconds()) % 60),
		"",
		p.loadPhaseMsg())

	return s
}

func (p *PomodoroModel) Quit() { 
	p.save()
}


func (p *PomodoroModel) GetPausedTime() time.Duration {
	return p.pausedTime
}

func (p *PomodoroModel) GetN() uint8 {
	return p.n
}

func (p *PomodoroModel) GetProgress() float64 {
	return float64(p.remainingTime) /
		   float64(p.phasesDurations[p.phaseType])
}

