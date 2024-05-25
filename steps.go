package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Steps struct {
	steps    []Step
	ready    bool
	viewport viewport.Model
	content  string
}

func initialSteps() Steps {
	type Config map[string][]string

	var ( //nolint:prealloc
		conf  Config
		steps []Step
	)

	dat, _ := os.ReadFile("./loci.toml")
	_, err := toml.Decode(string(dat), &conf)

	if err != nil {
		fmt.Println("loci.toml error", err) //nolint:forbidigo
	}

	for i, step := range conf["steps"] {
		steps = append(steps, initialStep(step, i))
	}

	return Steps{
		content:  "",
		steps:    steps,
		ready:    false,
		viewport: viewport.New(1, 1),
	}
}

func (m Steps) Init() tea.Cmd {
	var cmds []tea.Cmd //nolint:prealloc

	for _, s := range m.steps {
		cmds = append(cmds, s.Init())
	}

	cmd := func() tea.Msg { return startMsg{id: 0} }
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m Steps) UpdateAll(msg tea.Msg) (Steps, tea.Cmd) {
	var ( //nolint:prealloc
		cmd  tea.Cmd
		cmds []tea.Cmd
		skip bool
	)

	for index, s := range m.steps {
		currentStep, cmd := s.Update(msg)
		cmds = append(cmds, cmd)

		if skip {
			currentStep.state = Skipped
		}

		if currentStep.state == Exited1 {
			skip = true
		}

		m.steps[index] = currentStep
	}

	m.content = m.ViewAll()
	m.viewport.SetContent(m.content)
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Steps) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	nextModel, cmd := m.UpdateAll(msg)
	cmds = append(cmds, cmd)
	m = nextModel

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-lipgloss.Height(m.footerView()))
			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - lipgloss.Height(m.footerView())
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "r":
			for i := range m.steps {
				m.steps[i] = initialStep(m.steps[i].command, i)
			}

			return m, func() tea.Msg { return startMsg{id: 0} }
		}
	case exitMsg:
		nextID := msg.id + 1
		if nextID != len(m.steps) && msg.state != Exited1 {
			cmd := func() tea.Msg { return startMsg{id: msg.id + 1} }
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Steps) footerView() string {
	style := lipgloss.NewStyle().Background(lipgloss.Color("#bada55"))
	info := fmt.Sprintf("(r restart, q quit) %3.f%%", m.viewport.ScrollPercent()*100) //nolint:mnd,gomnd
	middle := strings.Repeat(" ", max(0, m.viewport.Width-lipgloss.Width(info)))

	return style.Render(lipgloss.JoinHorizontal(lipgloss.Center, middle, info))
}

func (m Steps) ViewAll() string {
	var (
		total   time.Duration
		content strings.Builder
	)

	for _, s := range m.steps {
		total += s.duration
		content.WriteString(s.View())
	}

	content.WriteString("⏱️  " + total.Round(time.Second).String())

	return content.String()
}

func (m Steps) View() string {
	return fmt.Sprintf("%s\n%s", m.viewport.View(), m.footerView())
}
