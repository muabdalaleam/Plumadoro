package main

import (
	"github.com/charmbracelet/lipgloss"
	// tea "github.com/charmbracelet/bubbletea"
)

const (
	errColor     = lipgloss.Color("1") // termenv's red
	warningColor = lipgloss.Color("3") // termenv's yellow
	alarmColor   = lipgloss.Color("4") // termenv's blue
	textColor    = lipgloss.Color("0") // termenv's black
)


func GetBorderStyle(color string) lipgloss.Style {
	var borderStyle lipgloss.Style
	var borderType  lipgloss.Border

	switch (Config.ProgressBar.Border) {
	case "rounded":   borderType = lipgloss.RoundedBorder()
	case "ascii":     borderType = lipgloss.ASCIIBorder()
	case "thick":     borderType = lipgloss.ThickBorder()
	case "double":    borderType = lipgloss.DoubleBorder()
	case "normal":    borderType = lipgloss.NormalBorder()
	case "hidden":    borderType = lipgloss.HiddenBorder()
	default:          borderType = lipgloss.NormalBorder()
	}

	borderStyle = lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(borderType).
		BorderForeground(lipgloss.Color(color))

	return borderStyle
}

func GetErrorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(0, 2).
		Background(lipgloss.Color(errColor)).
		Foreground(lipgloss.Color(textColor)).
		Width(int(Config.ProgressBar.MaxWidth))	
}

func GetWarningStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(0, 2).
		Background(lipgloss.Color(warningColor)).
		Foreground(lipgloss.Color(textColor)).
		Width(int(Config.ProgressBar.MaxWidth))	
}

func GetAlarmStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(0, 2).
		Background(lipgloss.Color(alarmColor)).
		Foreground(lipgloss.Color(textColor)).
		Width(int(Config.ProgressBar.MaxWidth))	
}


func GetCenterStyle(s string, height uint, width uint) lipgloss.Style {
	sw, sh := lipgloss.Width(s), lipgloss.Height(s)

	centerStyle := lipgloss.NewStyle().
		MarginLeft(max(0, (int(width)  - sw)/2)).
		MarginTop( max(0, (int(height) - sh)/2))

	return centerStyle
}
