package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type popupType int

type PopupMsg struct {
	Type    popupType
	Content string
}

type PopupModel struct {
	popup   PopupMsg
}

const (
	ErrorPopup popupType = iota
	WarningPopup
	AlarmPopup // XXX: rename this + find a use for it
)

func (m *PopupModel) Init() tea.Cmd {
	return nil
}

func (m *PopupModel) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			cmd = func() tea.Msg { return InitPomodoroMsg{} }
		}

	case PopupMsg:
		m.popup = msg
	}

	return cmd
}

func (m *PopupModel) Render() string {
	var popupStr string

	switch m.popup.Type { 
	case ErrorPopup:    popupStr = GetErrorStyle()  .Render(m.popup.Content, "\n")
	case WarningPopup:  popupStr = GetWarningStyle().Render(m.popup.Content, "\n")
	case AlarmPopup:    popupStr = GetAlarmStyle()  .Render(m.popup.Content, "\n")
	}

	content := lipgloss.JoinVertical(lipgloss.Center, "Ooops.", "")
	if len(popupStr) > 0 {
		content = lipgloss.JoinVertical(
			lipgloss.Center,
			content,
			lipgloss.JoinVertical(lipgloss.Center, popupStr),
		)
	}

	s := GetBorderStyle(Config.ProgressBar.PauseColor).Render(content)

	return s
}
