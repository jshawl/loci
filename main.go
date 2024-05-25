// package main

// import (
// 	"errors"
// 	"fmt"
// 	"os"
// 	"os/exec"
// 	"strings"
// 	"time"

// 	"github.com/BurntSushi/toml"
// )

// type Step struct {
// 	command  string
// 	duration time.Duration
// 	ok       bool
// }

// func (step Step) run() (Step, error) {
// 	start := time.Now()
// 	command := strings.Split(step.command, " ")
// 	cmd := exec.Command(command[0], command[1:]...)
// 	output, err := cmd.Output()
// 	step.duration = time.Since(start).Round(time.Millisecond)
// 	if err != nil {
// 		step.ok = false
// 		return step, errors.New(string(output))
// 	}
// 	step.ok = true
// 	return step, nil
// }

// func main() {
// 	type Config map[string][]string

// 	dat, _ := os.ReadFile("./loci.toml")

// 	var conf Config
// 	_, err := toml.Decode(string(dat), &conf)
// 	if err != nil {
// 		fmt.Println("loci.toml error", err)
// 	}

// 	var steps []Step

// 	for _, step := range conf["steps"] {
// 		steps = append(steps, Step{command: step})
// 	}

// 	for _, step := range steps {
// 		fmt.Print("running ", step.command, " ")
// 		step, err := step.run()
// 		var content strings.Builder
// 		if step.ok {
// 			content.WriteString("✅ ")
// 		} else {
// 			content.WriteString("❌ ")
// 		}
// 		content.WriteString(step.duration.String())
// 		content.WriteString("\n")
// 		fmt.Print(content.String())
// 		if err != nil {
// 			fmt.Println(err)
// 			break
// 		}
// 	}
// }

package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	steps    []Step
	quitting bool
}

func initialModel() model {
	return model{steps: []Step{
		newStep("go mod tidy"),
		newStep("go test -v ./..."),
	}}
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, s := range m.steps {
		cmds = append(cmds, s.Init())
	}
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
	}
	var cmds []tea.Cmd
	steps := m.steps
	m.steps = []Step{}
	for _, s := range steps {
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
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
