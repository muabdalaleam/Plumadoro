package main

import (
	"fmt"
	"math"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Kickstarts other necessary commands for the pomodoro model (tickers and window resizer)
type InitPomodoroMsg struct{}

type PomodoroTickMsg time.Time

type phaseType byte

type PomodoroModel struct {
	remainingTime    time.Duration
	pausedTime       time.Duration
	phaseType        phaseType
	running          bool
	n                uint8 

	// Configurable
	phasesDurations  map[phaseType]time.Duration
	progressBar      progress.Model
}

const (
	Focus phaseType = iota
	ShortBreak
	LongBreak
)


func (m *PomodoroModel) getPhaseMsg() string {
	var msg string

	if m.running {
		switch (m.phaseType) {
		case Focus:        msg = Config.ProgressBar.FocusMsg
		case ShortBreak:   msg = Config.ProgressBar.ShortBreakMsg
		case LongBreak:    msg = Config.ProgressBar.LongBreakMsg
		}
	} else {
		msg = Config.ProgressBar.PauseMsg
	}

	return msg
}

// HACK: IDK WTH is this
func (m *PomodoroModel) getPhaseColor() string { 
	var progressColor string
	switch (m.phaseType) {
		case Focus:       progressColor = Config.ProgressBar.FocusColor
		case ShortBreak:  progressColor = Config.ProgressBar.ShortBreakColor
		case LongBreak:   progressColor = Config.ProgressBar.LongBreakColor
	}

	return progressColor
}

func (m *PomodoroModel) getProgress() float64 {
	return float64(m.remainingTime) /
		   float64(m.phasesDurations[m.phaseType])
}


func (m *PomodoroModel) Init() tea.Cmd {
	if err := m.restore(); err != nil {
		*m = PomodoroModel {
			remainingTime: Config.Durations.Focus,
			pausedTime:    time.Duration(0),
			phaseType:     Focus,
			running:       Config.Autostart,
			n:             1, // NOTE: the index of phases is one based

			phasesDurations: map[phaseType]time.Duration{
				Focus:      Config.Durations.Focus,
				ShortBreak: Config.Durations.ShortBreak,
				LongBreak:  Config.Durations.LongBreak,
			},
			progressBar:     progress.New(
				progress.WithSolidFill(Config.ProgressBar.FocusColor)),
		}

		return func() tea.Msg { return PopupMsg{Type: ErrorPopup, Content: err.Error() } }
	}

	return func() tea.Msg { return InitPomodoroMsg{} } 
}

func (m *PomodoroModel) Update(msg tea.Msg) tea.Cmd { 
	var cmd tea.Cmd = nil

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// TODO: make the key bindings customizable
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.save()
			cmd = tea.Quit

		case " ":
			if !Config.Pausing && m.running {
				cmd = func() tea.Msg { return PopupMsg{
					Type: WarningPopup,
					Content: "Pausing phases is unallowed in your config",
				}}
			} else {
				m.toggle()
			}

		case "ctrl+r":
			m.reset()

		case "ctrl+s":
			if !Config.Skipping {
				cmd = func() tea.Msg { return PopupMsg{
					Type: WarningPopup,
					Content: "Skipping phases is unallowed in your config",
				}}
			} else {
				m.next()
			}
	}

	case tea.InterruptMsg, tea.QuitMsg:
		m.save()
		cmd = tea.Quit

	case tea.WindowSizeMsg:
		m.resizeProgressBar(msg.Width)

	case PomodoroTickMsg:
		m.tick(Config.TickDuration)
		cmd = tickPomodoroEvery()

		if !m.running {
			if m.pausedTime >= Config.MaxPauseDuration {
				cmd = func() tea.Msg { return PopupMsg{
					Type: WarningPopup,
					Content: "You have passed your maximum pause time per phase, resetting the phase.",
				}}
				m.reset()
			}
		}

	case InitPomodoroMsg:
		cmd = tea.Batch(
			tickPomodoroEvery(),
			tickLogEvery(),
			tea.WindowSize(),
		)

	case LogTickMsg:
		err := m.save()
		if err != nil {
			cmd = func() tea.Msg { return PopupMsg{Type: ErrorPopup, Content: err.Error()} }
		} else {
			cmd = tickLogEvery()
		}
	}

	return cmd
}

func (m *PomodoroModel) Render() string {
	phaseColor := m.getPhaseColor()
	if !m.running {
		phaseColor = Config.ProgressBar.PauseColor
	}

	s := GetBorderStyle(phaseColor).Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.NewStyle().Foreground(lipgloss.Color(phaseColor)).Render(m.getPhaseMsg()), 
			// make this configurable ^
			"",
			m.progressBar.ViewAs(m.getProgress()),
			"",
			fmt.Sprintf("Remaining: %02d:%02d | #%d",
				int(m.remainingTime.Minutes()),
				int(m.remainingTime.Seconds()) % 60,
				int(math.Ceil(float64(m.n) / 2.0))),
			),
		)

	return s
}


func (m *PomodoroModel) toggle() {
	if m.running == false {
		m.running = true
	} else {
		m.running = false
	}
}

func (m *PomodoroModel) tick(d time.Duration) {
	if m.running {
		m.remainingTime -= d
	} else {
		m.pausedTime += d
	}

	if m.remainingTime <= time.Duration(0) {
		m.next()
	}
}

func (m *PomodoroModel) reset() {
	m.pausedTime    = time.Duration(0)
	m.remainingTime = m.phasesDurations[m.phaseType]
	m.running       = Config.Autostart
	m.progressBar.FullColor = m.getPhaseColor()
}

// It updates the whole state of the PomodoroModel
func (m *PomodoroModel) next() {
	PlayAlarm() // HACK: i know this function shouldn't hanle alarms but u know

	// NOTE: be careful n is updated first
	m.n += 1

	var newPhaseType phaseType

	switch m.phaseType {
	case ShortBreak, LongBreak:
		newPhaseType = Focus
	// TODO: Checking next phase should be more configurable and not hard coded
	case Focus:
		if m.n == 7 || (m.n - 1) % 7 == 0 { // because for every 4 focus phases there are three short breaks in between
			newPhaseType = LongBreak
		} else {
			newPhaseType = ShortBreak
		}
	}

	m.pausedTime    = time.Duration(0)
	m.remainingTime = m.phasesDurations[newPhaseType]
	m.phaseType     = newPhaseType
	m.running       = Config.Autostart
	m.progressBar.FullColor = m.getPhaseColor()
}

func (m *PomodoroModel) resizeProgressBar(width int) {
	m.progressBar.Width = width - int(Config.ProgressBar.Padding) * 2 - 4
	if m.progressBar.Width > int(Config.ProgressBar.MaxWidth) {
		m.progressBar.Width = int(Config.ProgressBar.MaxWidth) 
	}
}

