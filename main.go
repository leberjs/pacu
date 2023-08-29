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
    err error
    profiles []string
    credentials
}

type credentials struct {
    key_id string
    secret_key string
    session_token string
}

type errMsg struct{ err error }

type profileMsg []string

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
    p := []string{}
    for _, ln := range lns {
        if strings.HasPrefix(ln, "[") {
            pn := strings.TrimSpace(ln[1:len(ln)-1])
            p = append(p, pn)
        }
    }

    return profileMsg(p)
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
        return m, tea.Quit
    
    case tea.KeyMsg:
        if msg.Type == tea.KeyCtrlC {
            return m, tea.Quit
        }
    }

    return m, nil
}

func (m model) View() string {
    if m.err != nil {
        return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
    }

    s := fmt.Sprintln("Grabbing profiles ... ")
    
    for _, p := range m.profiles {
        s += fmt.Sprintf("\nProfile Name: %s", p)
    }

    return "\n" + s + "\n\n"

}

func main() {
    p := tea.NewProgram(model{})
    if _, err := p.Run(); err != nil {
        log.Fatal(err)
    }
}
