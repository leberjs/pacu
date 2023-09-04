package internal

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const (
	textColor = lipgloss.Color("#007cbe")
)

var (
	inputStyle = lipgloss.NewStyle().Foreground(textColor)
)

func choicesView(m Model) string {
	var s string

	s += fmt.Sprintf("\n%s:\n\n", lipgloss.NewStyle().Bold(true).Render("Select Profile"))

	for i, p := range m.profiles {
		cursor := " "
		if m.cursor == i {
			cursor = "*"
		}
		v := inputStyle.Width(30).Render(p.name)
		s += fmt.Sprintf("%s  %s\n", cursor, v)
	}

	return s
}

func credentialEditView(m Model) string {
	var s string

	id := inputStyle.Width(6).Align(lipgloss.Left).Render("Key Id")
	secret := inputStyle.Width(10).Align(lipgloss.Left).Render("Secret Key")
	session := inputStyle.Width(13).Align(lipgloss.Left).Render("Session Token")

	s += fmt.Sprintf("\n%s:\n\n", lipgloss.NewStyle().Bold(true).Render("Enter Credentials"))

	s += fmt.Sprintf("%s: %s\n", id, m.InputState.inputs[0].View())
	s += fmt.Sprintf("%s: %s\n", secret, m.InputState.inputs[1].View())
	s += fmt.Sprintf("%s: %s\n", session, m.InputState.inputs[2].View())

	return s
}
