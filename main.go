package main

import (
	"fmt"
	"os"
	"strings"
	// "os"

	// "github.com/charmbracelet/lipgloss"
	// "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var choices = [...]string{"Tea", "Coffee", "Espresso", "Milk"}

type model struct {
	cursor int
	choice string
}

func (m model) Init() tea.Cmd {
	return nil
} 

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd = nil

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			cmd = tea.Quit

		case "j", "down":
			m.cursor++
			if m.cursor >= len(choices) {
				m.cursor = 0
			}

		case "k" ,"up":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(choices) - 1
			}
		}
	}

	m.choice = choices[m.cursor]

	return m, cmd
}

func (m model) View() string {
	s := strings.Builder{}

	s.WriteString("What kind of Bubble Tea would you like to order?\n\n")
	
	for i := 0; i < len(choices); i++ {
		if m.cursor == i {
			s.WriteString("(x)")
		} else {
			s.WriteString("( )")
		}

		s.WriteString(choices[i])
		s.WriteString("\n")
	}

	s.WriteString("\n(press q to quit)\n")

	return s.String()
}


func main() {
	p := tea.NewProgram(model{})

	// Run returns the model as a tea.Model.
	m, err := p.Run()
	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}

	// Assert the final tea.Model to our local model and print the choice.
	if m, ok := m.(model); ok && m.choice != "" {
		fmt.Printf("\n---\nYou chosed %s!\n\n", m.choice)
	}
}

