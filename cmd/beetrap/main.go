package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hnpf/beetrap/internal/capture"
	"github.com/hnpf/beetrap/internal/ui"
)

func main() {
	ch := make(chan capture.VxConnection, 32)
	p := tea.NewProgram(ui.VxNewModel(ch), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}