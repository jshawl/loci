package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type Step struct {
	command  string
	duration time.Duration
	ok       bool
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

func main() {
	type Config map[string][]string

	dat, _ := os.ReadFile("./loci.toml")

	var conf Config
	_, err := toml.Decode(string(dat), &conf)
	if err != nil {
		fmt.Println("loci.toml error", err)
	}

	var steps []Step

	for _, step := range conf["steps"] {
		steps = append(steps, Step{command: step})
	}

	for _, step := range steps {
		fmt.Print("running ", step.command, " ")
		step, err := step.run()
		var content strings.Builder
		if step.ok {
			content.WriteString("✅ ")
		} else {
			content.WriteString("❌ ")
		}
		content.WriteString(step.duration.String())
		content.WriteString("\n")
		fmt.Print(content.String())
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
