package main

import (
	"fmt"
	"os"

	"atlas.horizon/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

var Version = "dev"

func main() {
	// Version handling
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-v", "--version":
			fmt.Printf("atlas.horizon v%s\n", Version)
			return
		case "-h", "--help":
			fmt.Println("atlas.horizon - High-fidelity environmental and weather dashboard")
			fmt.Println("\nUsage:")
			fmt.Println("  atlas.horizon [flags]")
			fmt.Println("\nFlags:")
			fmt.Println("  -v, --version  Show version info")
			fmt.Println("  -h, --help     Show help info")
			return
		}
	}

	p := tea.NewProgram(ui.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
