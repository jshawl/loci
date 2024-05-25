package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Step struct {
	command  string
	spinner  spinner.Model
	duration time.Duration
	ok       bool
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

func (step Step) run() (Step, error) {
	start := time.Now()
	command := strings.Split(step.command, " ")
	cmd := exec.Command(command[0], command[1:]...)
	output, err := cmd.Output()
	step.duration = time.Since(start).Round(time.Millisecond)
	if err != nil {
		step.ok = false
		return step, errors.New(string(output))
	}
	step.ok = true
	return step, nil
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
