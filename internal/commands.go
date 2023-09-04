package internal

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) getProfilesCmd() tea.Msg {
	bs, err := os.ReadFile(m.credentialsPath)
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

	return profileFetchMsg(ps)
}

func (m *Model) updateInputTextCmd(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *Model) setNewCredentialsCmd() tea.Cmd {
	return func() tea.Msg {
		p := m.profiles[m.ProfileState.selectedIndex]

		p.credentials.key_id = m.inputs[0].Value()
		p.credentials.secret_key = m.inputs[1].Value()
		p.credentials.session_token = m.inputs[2].Value()

		m.profiles[m.ProfileState.selectedIndex] = p

		return profileSelectedMsg{}
	}
}

func writeProfilesCmd(m Model) tea.Cmd {
	return func() tea.Msg {
		f, err := os.Create(m.credentialsPath)
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

		return credentialsFileWrittenMsg{}
	}
}
