package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	cursor int
	credentials
	err      error
	profiles []profile
	selected int
}

type profile struct {
	name string
	credentials
}

type credentials struct {
	key_id        string
	secret_key    string
	session_token string
}

type errMsg struct{ err error }

type profileMsg []profile

const credentialPath = ".aws/credentials"

func (e errMsg) Error() string { return e.err.Error() }

func getProfiles() tea.Msg {
	dn, err := os.UserHomeDir()
	path := filepath.Join(dn, credentialPath)

	bs, err := os.ReadFile(path)
	if err != nil {
		return errMsg{err}
	}

	lns := strings.Split(string(bs), "\n")
	ps := []profile{}
	for idx, ln := range lns {
		if strings.HasPrefix(ln, "[") {
			p := profile{}
			pn := strings.TrimSpace(ln[1 : len(ln)-1])
			p.name = pn

			cred := lns[idx+1 : idx+4]

			for _, cl := range cred {
				if strings.HasPrefix(cl, "aws_access_key_id") {
					p.credentials.key_id = assignCred(cl)
				} else if strings.HasPrefix(cl, "aws_secret_access_key") {
					p.credentials.key_id = assignCred(cl)
				} else if strings.HasPrefix(cl, "aws_session_token") {
					p.credentials.key_id = assignCred(cl)
				}
			}
			ps = append(ps, p)
		}
	}

	return profileMsg(ps)
}

func assignCred(s string) string {
	ss := strings.Split(s, "=")
	return strings.TrimSpace(ss[1])
}

func (m model) Init() tea.Cmd {
	return getProfiles
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case errMsg:
		m.err = msg
		return m, tea.Quit

	case profileMsg:
		m.profiles = msg

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.profiles)-1 {
				m.cursor++
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}

	s := ""

	if m.selected > -1 {
		s = "I selected something!\n"
	}

	if len(m.profiles) == 0 {
		s = fmt.Sprintln("Grabbing profiles ... ")
	} else {
		s = ""
		for i, p := range m.profiles {
            cursor := " "
			if m.cursor == i {
			    cursor = "*"
			}
			s += fmt.Sprintf("%s  %s\n", cursor, p.name)
		}
	}

	s += "\nCtrl+C to quit\n"

	return s

}

func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
