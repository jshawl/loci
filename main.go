package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/BurntSushi/toml"
)

func main() {
	type Config map[string][]string

	dat, _ := os.ReadFile("./loci.toml")

	var conf Config
	_, err := toml.Decode(string(dat), &conf)
	if err != nil {
		fmt.Println("loci.toml error", err)
	}

	for _, step := range conf["steps"] {
		command := strings.Split(step, " ")
		fmt.Println("running", step, "...")
		cmd := exec.Command(command[0], command[1:]...)
		output, err := cmd.Output()
		if err != nil {
			fmt.Println(string(output))
		}
	}
}
