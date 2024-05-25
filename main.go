package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	steps    []Step
	quitting bool
}

// for _, step := range steps {
// 	fmt.Print("running ", step.command, " ")
// 	step, err := step.run()
// 	var content strings.Builder
// 	if step.ok {
// 		content.WriteString("✅ ")
// 	} else {
// 		content.WriteString("❌ ")
// 	}
// 	content.WriteString(step.duration.String())
// 	content.WriteString("\n")
// 	fmt.Print(content.String())
// 	if err != nil {
// 		fmt.Println(err)
// 		break
// 	}
// }

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
	cmds = append(cmds, m.steps[0].start())
	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}
	case exitMsg:
		log.Println("received an exitMsg in main.go", msg)
	}
	var cmds []tea.Cmd
	steps := m.steps
	m.steps = []Step{}
	for _, s := range steps {
		// run msg will go through here
		s, cmd := s.Update(msg)
		m.steps = append(m.steps, s)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var content strings.Builder
	for _, s := range m.steps {
		content.WriteString(s.View())
	}
	str := content.String()
	if m.quitting {
		return str + "\n"
	}
	return str
}

func main() {
	p := tea.NewProgram(initialModel())
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
