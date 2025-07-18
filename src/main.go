// With the name of the Lord the most merciful, most knowledgable.

package main

import (
	"fmt"
	"os"
	"plumadoro/config"
	"plumadoro/pomodoro"
	"strings"
	"time"

	// "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type styleID byte

type PlumaModel struct {
	pomodoro    *pomodoro.PomodoroModel

	height      int  // HACK: i think uint16 is more suitable
	width       int
	styles      map[styleID]lipgloss.Style

	errors      []error // errors in this are the non-terminating errors such as failed config loading
}

const (
	errorStyle styleID = iota
)


func tickEvery(dur time.Duration) tea.Cmd {
	return tea.Every(dur, func(t time.Time) tea.Msg {
		return pomodoro.TickMsg(t)
	})
}

func (pm *PlumaModel) Init() tea.Cmd {
	// Loading the config
	err := config.LoadConfig()

	if err != nil {
		pm.errors = append(pm.errors, err)
	}

	pm.pomodoro = &pomodoro.PomodoroModel{}
	pm.pomodoro.Init()

	// TODO: load those styles options from the config
	pm.styles = map[styleID]lipgloss.Style{}

	pm.styles[errorStyle] = lipgloss.NewStyle().
		Padding(0, 2).
		Background(lipgloss.Color("#FF0000")).
		Foreground(lipgloss.Color("#FFFFFF"))

	return tea.Batch(
		tickEvery(config.Config.TickDuration),
		tea.WindowSize(),
	)
}

func (pm *PlumaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			pm.pomodoro.Quit()
			cmd = tea.Quit

		case " ":
			pm.pomodoro.Toggle()
	}

	case tea.InterruptMsg, tea.QuitMsg:
		pm.pomodoro.Quit()
		cmd = tea.Quit

	case pomodoro.TickMsg:
		pm.pomodoro.Tick(config.Config.TickDuration)
		cmd = tickEvery(config.Config.TickDuration)

	case tea.WindowSizeMsg:
		pm.width  = msg.Width
		pm.height = msg.Height

		pm.pomodoro.ProgressBar.Width = msg.Width - int(config.Config.ProgressBar.Padding) * 2 - 4
		if pm.pomodoro.ProgressBar.Width > int(config.Config.ProgressBar.MaxWidth) {
			pm.pomodoro.ProgressBar.Width = int(config.Config.ProgressBar.MaxWidth) 
		}
	}

	return pm, cmd
}

func (pm *PlumaModel) View() string {
	var s string 

	// Viewing the progressBar with a Remaining time bar
	s = pm.pomodoro.Render()

	// Centering the view
	sw, sh := lipgloss.Width(s), lipgloss.Height(s)
	centerStyle := lipgloss.NewStyle().
		MarginLeft(max(0, (pm.width-sw)/2)).
		MarginTop( max(0, (pm.height-sh)/2))
	
	s = centerStyle.Render(s)

	// Viewing the errors.
	if len(pm.errors) != 0 {
		var errMsg strings.Builder
		for _, err := range pm.errors {
			errMsg.WriteString(err.Error())
		}

		s = lipgloss.JoinVertical(
			lipgloss.Left,
			s,
			pm.styles[errorStyle].MarginTop(max(0, (pm.height-sh)/2)).
				Render(errMsg.String()),
		)
	}

	return s
}


func main() {
	// TODO: support command line arguments
	p := tea.NewProgram(&PlumaModel{},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err := p.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

