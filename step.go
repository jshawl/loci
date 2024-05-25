package main

import (
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
	id       int
	state    StepState
	output   string
	duration time.Duration
}

type StepState string

const (
	Pending StepState = "🔜"
	Started StepState = ""
	Exited0 StepState = "✅"
	Exited1 StepState = "❌"
	Skipped StepState = "⏭️"
)

type Step struct {
	id        int
	command   string
	spinner   spinner.Model
	startedAt time.Time
	duration  time.Duration
	state     StepState
	output    string
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

func (m Step) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Step) start() (Step, tea.Cmd) {
	log.Println("starting...", m.command)
	m.state = Started
	m.startedAt = time.Now()
	return m, func() tea.Msg {
		start := time.Now()
		command := strings.Split(m.command, " ")
		cmd := exec.Command(command[0], command[1:]...)
		output, err := cmd.Output()
		m.duration = time.Since(start).Round(time.Millisecond)
		if err != nil {
			m.state = Exited1
		} else {
			m.state = Exited0
		}
		return exitMsg{
			id:       m.id,
			state:    m.state,
			output:   string(output),
			duration: m.duration,
		}
	}
}

func (m Step) Update(msg tea.Msg) (Step, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case startMsg:
		if m.id == msg.id {
			m, cmd := m.start()
			return m, cmd
		}
	case exitMsg:
		if m.id == msg.id {
			m.state = msg.state
			m.duration = msg.duration
			m.output = msg.output
		}
	case spinner.TickMsg:
		if m.state == Started || m.state == Pending {
			m.spinner, cmd = m.spinner.Update(msg)
		}
		return m, cmd
	}
	if m.state == Started {
		m.duration = time.Since(m.startedAt).Round(time.Millisecond)
	}
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
		return fmt.Sprintf("%s %s %s\n %s", icon, m.command, m.duration, m.output)
	}
	if m.state == Skipped {
		icon = string(Skipped)
		return fmt.Sprintf("%s  %s (skipped)\n", icon, m.command)
	}
	if m.state != Pending {
		return fmt.Sprintf("%s %s %s\n", icon, m.command, m.duration)
	}
	return fmt.Sprintf("%s %s \n", icon, m.command)
}
