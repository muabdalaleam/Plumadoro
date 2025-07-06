// With the name of the Lord the most merciful, most knowledgable.

package main

import (
	"fmt"
	"os"
	"plumadoro/phase"
	"time"
	// "strings"

	// "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/progress"
)

type tickMsg time.Time

type PlumaModel struct {
	phase  phase.Phase

	progressBar progress.Model
	height      uint16
	width       uint16
}

// TODO: make those values configrable
const (
	tickDuration = time.Millisecond * 20
	padding      = 4
	maxWidth     = 80
)

var phaseDurations = map[phase.PhaseType]time.Duration {
	phase.Focus:      time.Minute * 20,
	phase.ShortBreak: time.Minute * 5,
	phase.LongBreak:  time.Minute * 20,
}


// TODO: make tick speed configurable
func tickEvery() tea.Cmd {
	return tea.Every(tickDuration, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (pm *PlumaModel) Init() tea.Cmd {
	phase.SetPhaseDurations(phaseDurations)

	pm.phase = phase.InitFirstPhase()
	pm.progressBar = progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))

	// HACK: kickstarting the window size by sending a WindowSizeMsg isn't ideal
	return tea.Batch(
		tickEvery(),
		tea.WindowSize(),
	)}

func (pm *PlumaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			cmd = tea.Quit

		case " ":
			pm.phase.Toggle()
	}

	case tea.InterruptMsg, tea.QuitMsg:
		cmd = tea.Quit

	case tickMsg:
		pm.phase.Tick(tickDuration)
		cmd = tickEvery()

	case progress.FrameMsg:
		var progressModel tea.Model

		progressModel, cmd = pm.progressBar.Update(msg)
		pm.progressBar = progressModel.(progress.Model)

	case tea.WindowSizeMsg:
		pm.width  = uint16(msg.Width)
		pm.height = uint16(msg.Height)

		pm.progressBar.Width = msg.Width - padding * 2 - 4
		if pm.progressBar.Width > maxWidth {
			pm.progressBar.Width = maxWidth
		}
	}

	return pm, cmd
}

func (pm *PlumaModel) View() string {
	wrapperStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#FFFFFF"))

	s := lipgloss.JoinVertical(
		lipgloss.Center,
		wrapperStyle.Render(pm.progressBar.ViewAs(pm.phase.GetProgress())),
		"",
		fmt.Sprintf("Remaining: %02d:%02d",
			int32(pm.phase.GetRemainingTime().Minutes()),
			int32(pm.phase.GetRemainingTime().Seconds()) % 60),
	)

	// Kudos for GPT 4.o mini
	w, h := int(pm.width), int(pm.height)
	sw, sh := lipgloss.Width(s), lipgloss.Height(s)

	centerStyle := lipgloss.NewStyle().
		MarginLeft(max(0, (w-sw)/2)).
		MarginTop(max(0, (h-sh)/2))

	s = centerStyle.Render(s)
	return s
	// s.WriteString(pm.phase.GetRemainingTime().String())
}

func main() {
	// TODO: support command line arguments
	p := tea.NewProgram(&PlumaModel{},
		tea.WithAltScreen(), // TODO: alot of those options should be configrable
		tea.WithMouseCellMotion(),
		// tea.WithFPS(30),
	)

	_, err := p.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

