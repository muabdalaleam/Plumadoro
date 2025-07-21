// With the name of the Lord the most merciful, most knowledgable.

package main

import (
	"fmt"
	"os"
	"time"

	// "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type MainModel struct {
	pomodoro    *PomodoroModel
	popup       *PopupModel
	// help        *HelpModel

	height      int  // HACK: i think uint16 is more suitable
	width       int

	activeSubmodel Submodel
}


func tickLogEvery() tea.Cmd {
	return tea.Every(LogTickDuration, func(t time.Time) tea.Msg { return LogTickMsg(t) } )
}

func tickPomodoroEvery() tea.Cmd {
	return tea.Every(Config.TickDuration, func(t time.Time) tea.Msg { return PomodoroTickMsg(t) } )
}

func (m *MainModel) Init() tea.Cmd {
	var cmd tea.Cmd

	// Loading the config
	// Checking the config errors after the popup model has been intialized
	err := LoadConfig()
	m.popup    = &PopupModel{}
	m.pomodoro = &PomodoroModel{}
	cmd = m.pomodoro.Init()

	if err != nil {
		cmd = tea.Batch(
			cmd,
			func() tea.Msg { return PopupMsg{Type: ErrorPopup, Content: err.Error()} },
		)
	}

	return cmd
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil

	switch (m.activeSubmodel) {
	case m.pomodoro: cmd = m.pomodoro.Update(msg)
	case m.popup:    cmd = m.popup.Update(msg)
	}


	switch msg := msg.(type) {
	case InitPomodoroMsg:
		if m.activeSubmodel != m.pomodoro {
			cmd = tea.Batch(cmd, func() tea.Msg { return msg } )
		}
		m.activeSubmodel = m.pomodoro

	case PopupMsg:
		if m.activeSubmodel != m.popup {
			cmd = tea.Batch(cmd, func() tea.Msg { return msg })
		}
		m.activeSubmodel = m.popup

	case tea.InterruptMsg, tea.QuitMsg:
		cmd = tea.Batch(cmd, tea.Quit)

	case tea.WindowSizeMsg:
		m.width  = msg.Width
		m.height = msg.Height
	}

	return m, cmd
}

func (m *MainModel) View() string {
	var s string 

	// Viewing the progressBar with a Remaining time bar
	switch (m.activeSubmodel) {
	case m.pomodoro:  s = m.pomodoro.Render()
	case m.popup:     s = m.popup.Render()
	}
	
	// Centering the view
	s = GetCenterStyle(s, uint(m.height), uint(m.width)).
		Render(s)

	// Viewing the errors.
	// if len(m.errors) != 0 {
	// 	var errMsg strings.Builder
	// 	for _, err := range m.errors {
	// 		errMsg.WriteString(err.Error())
	// 	}

	// 	s = lipgloss.JoinVertical(
	// 		lipgloss.Left,
	// 		s,
	// 		m.styles[errorStyle].MarginTop(max(0, (m.height-sh)/2)).
	// 			Render(errMsg.String()),
	// 	)
	// }

	return s
}


func main() {
	// TODO: support command line arguments
	p := tea.NewProgram(&MainModel{},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err := p.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

