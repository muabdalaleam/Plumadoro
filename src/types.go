package main

import (
	"github.com/charmbracelet/bubbletea"
)


type Submodel interface {
	Init() tea.Cmd
	Update(msg tea.Msg) tea.Cmd
	Render() string
} 
