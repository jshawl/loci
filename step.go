package main

import (
	"fmt"
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
	Skipped StepState = "🙈"
)

type Step struct {
	command   string
	duration  time.Duration
	id        int
	output    string
	spinner   spinner.Model
	startedAt time.Time
	state     StepState
}

func initialStep(command string, index int) Step {
	stepSpinner := spinner.New()
	stepSpinner.Spinner = spinner.Line
	stepSpinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Step{
		command:   command,
		duration:  0,
		id:        index,
		output:    "",
		spinner:   stepSpinner,
		startedAt: time.Now(),
		state:     Pending,
	}
}

func (m Step) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Step) start() (Step, tea.Cmd) {
	m.state = Started
	m.startedAt = time.Now()

	return m, func() tea.Msg {
		start := time.Now()
		command := strings.Split(m.command, " ")
		cmd := exec.Command(command[0], command[1:]...) //nolint:gosec
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
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	if m.state == Started {
		m.duration = time.Since(m.startedAt).Round(time.Millisecond)
	}

	switch msg := msg.(type) {
	case startMsg:
		if m.id == msg.id {
			m, cmd := m.start()
			cmds = append(cmds, cmd, m.spinner.Tick)

			return m, tea.Batch(cmds...)
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

	return m, cmd
}

func (m Step) View() string {
	var icon string
	if m.state == Pending {
		icon = string(Pending)
	}

	if m.state == Started {
		icon = m.spinner.View() + " "
	}

	if m.state == Exited0 {
		icon = string(Exited0)
	}

	if m.state == Exited1 {
		icon = string(Exited1)
		style := lipgloss.NewStyle().
			Border(lipgloss.HiddenBorder(), false, false, false, true).
			BorderBackground(lipgloss.Color("#FF0000")).
			Padding(0, 0, 0, 1)

		return fmt.Sprintf("%s %s %s\n%s\n", icon, m.command, m.duration, style.Render(m.output))
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
