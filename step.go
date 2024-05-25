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

type startMsg struct {
	id int
}

type exitMsg struct {
	id    int
	state StepState
}

type StepState string

const (
	Pending StepState = "🔜"
	Started StepState = ""
	Exited0 StepState = "✅"
	Exited1 StepState = "❌"
)

type Step struct {
	id       int
	command  string
	spinner  spinner.Model
	duration time.Duration
	state    StepState
}

func newStep(command string, id int) Step {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return Step{
		command: command,
		spinner: s,
		id:      id,
		state:   Pending,
	}
}

func (step Step) run() (Step, error) {
	start := time.Now()
	command := strings.Split(step.command, " ")
	cmd := exec.Command(command[0], command[1:]...)
	output, err := cmd.Output()
	step.duration = time.Since(start).Round(time.Millisecond)
	if err != nil {
		step.state = Exited1
		return step, errors.New(string(output))
	}
	step.state = Exited0
	return step, nil
}

func (m Step) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Step) start() (Step, tea.Cmd) {
	log.Println("starting...")
	m.state = Started
	return m, func() tea.Msg {
		// i/o
		time.Sleep(time.Second)
		return exitMsg{id: m.id, state: Exited1}
	}
}

func (m Step) Update(msg tea.Msg) (Step, tea.Cmd) {
	switch msg := msg.(type) {
	case startMsg:
		if m.id == msg.id {
			log.Println("received a startMsg in step.go", msg)
			m, cmd := m.start()
			return m, cmd
		}
	case exitMsg:
		if m.id == msg.id {
			// m.ok = true
			m.state = msg.state
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

	var icon string

	if m.state == Pending {
		icon = string(Pending)
	}
	if m.state == Started {
		icon = fmt.Sprintf("%s ", m.spinner.View())
	}
	if m.state == Exited0 {
		icon = string(Exited0)
	}
	if m.state == Exited1 {
		icon = string(Exited1)
	}
	return fmt.Sprintf("%s %s\n", icon, m.command)
}
