package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Step struct {
	command string
	spinner spinner.Model
}

func newStep(command string) Step {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return Step{
		command: command,
		spinner: s,
	}
}

func (m Step) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Step) Update(msg tea.Msg) (Step, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m Step) View() string {
	return fmt.Sprintf("%s %s\n", m.spinner.View(), m.command)
}
