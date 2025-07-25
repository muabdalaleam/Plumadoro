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

type ResetPopupsMsg struct {}

type PopupModel struct {
	popups   []PopupMsg
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

		case "ctrl+r":
			m.popups = []PopupMsg{} // resetting
		}

	case PopupMsg:
		m.popups = append(m.popups, msg)
	}

	return cmd
}

func (m *PopupModel) Render() string {
	var popupsStr string

	for i, popup := range m.popups {
		newline := "\n\n"
		if i == len(m.popups) - 1 {
			newline = ""
		}
		switch popup.Type {
		case ErrorPopup:   popupsStr += GetErrorStyle().Render(popup.Content)   + newline
		case WarningPopup: popupsStr += GetWarningStyle().Render(popup.Content) + newline
		case AlarmPopup:   popupsStr += GetAlarmStyle().Render(popup.Content)   + newline
		}
	}


	content := lipgloss.JoinVertical(lipgloss.Center, "Ooops.", "")
	content = lipgloss.JoinVertical(
		lipgloss.Center,
		content,
		lipgloss.JoinVertical(lipgloss.Center, popupsStr),
	)

	s := GetBorderStyle(Config.ProgressBar.PauseColor).Render(content)

	return s
}
