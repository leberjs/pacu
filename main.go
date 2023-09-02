package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	cursor int
	credentials
	err        error
	focusIndex int
	inputs     []textinput.Model
	profiles   []profile
	selected   int
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

type profileUpdateMsg struct{}

type completeMsg struct{}

const credentialPath = ".aws/credentials"

func (e errMsg) Error() string { return e.err.Error() }

func initialModel() model {
	m := model{
		selected: -1,
		inputs:   make([]textinput.Model, 3),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CharLimit = 128

		switch i {
		case 0:
			t.Placeholder = "aws_access_key_id"
			t.Focus()
		case 1:
			t.Placeholder = "aws_secret_access_key"
		case 2:
			t.Placeholder = "aws_session_token"
		}

		m.inputs[i] = t
	}

	return m
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *model) updateSelectedProfile() {
	p := m.profiles[m.selected]

	p.credentials.key_id = m.inputs[0].Value()
	p.credentials.secret_key = m.inputs[1].Value()
	p.credentials.session_token = m.inputs[2].Value()

	m.profiles[m.selected] = p
}

func quitCli() tea.Msg {
    return completeMsg{}
}

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
					p.credentials.secret_key = assignCred(cl)
				} else if strings.HasPrefix(cl, "aws_session_token") {
					p.credentials.session_token = assignCred(cl)
				}
			}
			ps = append(ps, p)
		}
	}

	return profileMsg(ps)
}

func writeProfiles(m model) {
	dn, err := os.UserHomeDir()
	path := filepath.Join(dn, credentialPath)

	f, err := os.Create(path)
	if err != nil {
		m.err = err
		tea.Quit()
	} else {
		var s string
		for _, p := range m.profiles {
			s += "[" + p.name + "]\n"
			s += "aws_access_key_id = " + p.credentials.key_id + "\n"
			s += "aws_secret_access_key = " + p.credentials.secret_key + "\n"
			if p.credentials.session_token != "" {
				s += "aws_session_token = " + p.credentials.session_token + "\n"
			}
			s += "\n"
		}

		f.WriteString(s)
	}
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

	case completeMsg:
        return m, tea.Quit

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

		case "enter":
			if m.selected < 0 {
				m.selected = m.cursor
			} else {
				if m.focusIndex == len(m.inputs)-1 {
					m.updateSelectedProfile()
					writeProfiles(m)
                    return m, quitCli
				} else {
					m.inputs[m.focusIndex].Blur()
					m.focusIndex++
					m.inputs[m.focusIndex].Focus()
				}
			}
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}

	s := ""

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

	if m.selected > -1 {
		s = ""
		for i := range m.inputs {
			s += m.inputs[i].View()
			if i < len(m.inputs)-1 {
				s += "\n"
			}
		}
		s += "\n"
	}

	s += "\nCtrl+C to quit\n"

	return s

}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
