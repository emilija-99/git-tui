package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	p := tea.NewProgram(New(dir))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
