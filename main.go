package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	steps    []Step
	ready    bool
	quitting bool
	viewport viewport.Model
	content  string
}

func initialModel() model {
	type Config map[string][]string
	dat, _ := os.ReadFile("./loci.toml")
	var conf Config
	_, err := toml.Decode(string(dat), &conf)
	if err != nil {
		fmt.Println("loci.toml error", err)
	}
	var steps []Step
	for i, step := range conf["steps"] {
		steps = append(steps, newStep(step, i))
	}
	return model{steps: steps}
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, s := range m.steps {
		cmds = append(cmds, s.Init())
	}
	cmd := func() tea.Msg { return startMsg{id: 0} }
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	steps := m.steps
	m.steps = []Step{}
	var content strings.Builder
	var skip bool
	for _, s := range steps {
		s, cmd := s.Update(msg)
		m.steps = append(m.steps, s)
		cmds = append(cmds, cmd)
		if skip {
			s.state = Skipped
		}
		content.WriteString(s.View())
		if s.state == Exited1 {
			skip = true
		}
	}
	m.content = content.String()
	m.viewport.SetContent(m.content)
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height)
			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "r":
			steps := m.steps
			for i := range steps {
				m.steps[i].state = Pending
			}
			return m, func() tea.Msg { return startMsg{id: 0} }
		default:
			return m, nil
		}
	case exitMsg:
		log.Println("received an exitMsg in main.go", msg)
		cmd := func() tea.Msg { return startMsg{id: msg.id + 1} }
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return m.viewport.View()
}

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		f.Truncate(0)
		f.Seek(0, 0)
		log.Println("program starting...")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
