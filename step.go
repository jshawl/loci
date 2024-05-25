package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// sent by a step
type exitMsg struct {
	id int
}

type Step struct {
	id       int
	command  string
	spinner  spinner.Model
	duration time.Duration
	ok       bool
}

func newStep(command string, id int) Step {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return Step{
		command: command,
		spinner: s,
		id:      id,
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

func (m Step) start() tea.Cmd {
	log.Println("starting...")
	return func() tea.Msg {
		// i/o
		time.Sleep(time.Second)
		return exitMsg{id: m.id}
	}
}

func (m Step) Update(msg tea.Msg) (Step, tea.Cmd) {
	switch msg := msg.(type) {
	case exitMsg:
		if m.id == msg.id {
			m.ok = true
			log.Println("received an exitMsg in step.go", msg)
		}
	}
	// log.Println("step.Update", msg)
	// receive run message
	// -> append exitmsg after run to cmd
	// case startMsg:
	// startMsg.id == m.id
	// m, cmd := m.start()
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m Step) View() string {
	if m.ok {
		return "ok\n"
	} else {
		return fmt.Sprintf("%s %s\n", m.spinner.View(), m.command)
	}
}
