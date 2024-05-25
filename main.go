package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	program := tea.NewProgram(
		initialSteps(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if len(os.Getenv("DEBUG")) > 0 {
		file, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err) //nolint:forbidigo

			defer func() {
				os.Exit(1)
			}()
		}
		defer file.Close()

		_ = file.Truncate(0)
		_, _ = file.Seek(0, 0)
	}

	if _, err := program.Run(); err != nil {
		fmt.Println(err) //nolint:forbidigo

		defer func() {
			os.Exit(1)
		}()
	}
}
