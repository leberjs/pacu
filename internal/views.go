package internal

import (
	"fmt"
)

func choicesView(m Model) string {
	var s string

    s += "\nSelect Profile:\n\n"

	for i, p := range m.profiles {
		cursor := " "
		if m.cursor == i {
			cursor = "*"
		}
		s += fmt.Sprintf("%s  %s\n", cursor, p.name)
	}

	return s
}

func credentialEditView(m Model) string {
    var s string

    s += "\nEnter Credentials:\n\n"

	for i := range m.inputs {
		s += m.inputs[i].View()
        s += "\n"
	}

    return s
}
