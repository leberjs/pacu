package internal

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func Run() {
	p := tea.NewProgram(NewModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
