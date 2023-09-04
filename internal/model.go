package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	credentialsPath string
	cursor          int
	err             error
	InputState
	ProfileState
}

// Input State
type InputState struct {
	focusIndex int
	inputs     []textinput.Model
}

// Profile State
type ProfileState struct {
	profiles      []profile
	selectedIndex int
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

func NewModel() Model {
	dir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("\nWe had some trouble: %v\n\n", err)
	}

	m := Model{
		credentialsPath: filepath.Join(dir, ".aws/credentials"),
		InputState: InputState{
			inputs: make([]textinput.Model, 3),
		},
		ProfileState: ProfileState{
			selectedIndex: -1,
		},
	}

	var t textinput.Model
	for i := range m.InputState.inputs {
		t = textinput.New()
		t.Prompt = ""
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

func (m Model) Init() tea.Cmd {
	return m.getProfilesCmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.err = msg
		return m, tea.Quit

	case profileFetchMsg:
		m.profiles = msg

	case profileSelectedMsg:
		return m, writeProfilesCmd(m)

	case credentialsFileWrittenMsg:
		return m, tea.Quit

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}

	if m.ProfileState.selectedIndex == -1 {
		return updateProfileChoice(msg, m)
	} else {
		return updateCredentialEdit(msg, m)
	}
}

func updateProfileChoice(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Up):
			if m.cursor > 0 {
				m.cursor--
			}

		case key.Matches(msg, DefaultKeyMap.Down):
			if m.cursor < len(m.profiles)-1 {
				m.cursor++
			}

		case key.Matches(msg, DefaultKeyMap.Enter):
			m.ProfileState.selectedIndex = m.cursor
		}
	}

	return m, nil
}

func updateCredentialEdit(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Enter):
			if m.InputState.focusIndex == len(m.InputState.inputs)-1 {
				return m, m.setNewCredentialsCmd()
			} else {
				m.InputState.inputs[m.InputState.focusIndex].Blur()
				m.InputState.focusIndex++
				m.InputState.inputs[m.InputState.focusIndex].Focus()
			}

		}
	}

	cmd := m.updateInputTextCmd(msg)

	return m, cmd
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}

	var s string

	if m.ProfileState.selectedIndex == -1 {
		s = choicesView(m)
	} else {
		s = credentialEditView(m)
	}

	s += "\nCtrl+C to quit\n"

	return s
}
