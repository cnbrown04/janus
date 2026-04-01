package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	// Import your new UI package
	"github.com/cnbrown04/janus/ui" 
)

func main() {
	// Initialize your exported Model using ui.New()
	p := tea.NewProgram(ui.New(), tea.WithAltScreen())
	
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
