package ui

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mellonnen/chronograph/git"
	"github.com/mellonnen/chronograph/models"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	titleStyle          = lipgloss.NewStyle().Bold(true)

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type formModel struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode textinput.CursorMode

	keys     formKeyMap
	resource Resource
	title    string
}

type formKeyMap struct {
	next key.Binding
	prev key.Binding
}

func newFormKeyMap() formKeyMap {
	return formKeyMap{
		next: key.NewBinding(key.WithKeys("down", "tab", "enter")),
		prev: key.NewBinding(key.WithKeys("up", "shift+tab")),
	}
}

func newForm(r Resource) formModel {
	m := formModel{}
	m.resource = r
	m.keys = newFormKeyMap()
	m.title = strings.Title(fmt.Sprintf("create new %s", r))
	m.inputs = make([]textinput.Model, 0)
	m.inputs = append(m.inputs, createTextInput("Name"))
	m.inputs = append(m.inputs, createTextInput("Description"))
	m.inputs[0].Focus()
	m.inputs[0].PromptStyle = focusedStyle
	m.inputs[0].TextStyle = focusedStyle

	switch r {
	case Repo:
		pathInput := createTextInput("Path to Repo")
		cwd, _ := os.Getwd()
		path, _ := git.RepoPathFromPath(cwd)
		pathInput.SetValue(path)
		m.inputs = append(m.inputs, pathInput)

	}
	return m
}

func (f formModel) init() tea.Cmd {
	return textinput.Blink
}

func (m formModel) update(msg tea.Msg) (formModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {

		case key.Matches(msg, m.keys.prev), key.Matches(msg, m.keys.next):
			if msg.String() == "enter" && m.focusIndex == len(m.inputs) {

				switch m.resource {
				case Workspace:
					workspace := models.Workspace{
						Name:        m.inputs[0].Value(),
						Description: sql.NullString{String: m.inputs[1].Value(), Valid: true},
					}

					return m, addWorkspaceCmd(workspace)
				case Repo:
					path := m.inputs[2].Value()
					remote, err := git.RemoteFromPath(path)
					if err != nil {
						return m, errorCmd(fmt.Errorf("getting remote: %v", err))
					}
					repo := models.Repo{
						Name:        m.inputs[0].Value(),
						Description: sql.NullString{String: m.inputs[1].Value(), Valid: true},
						Path:        path,
						Remote:      remote,
					}
					return m, addRepoCmd(repo)
				}
			}

			if key.Matches(msg, m.keys.prev) {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}
			return m, tea.Batch(cmds...)
		}
	}
	// Update cursor blinking.
	var cmds = make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m formModel) view() string {
	var b strings.Builder

	fmt.Fprintf(&b, "%s\n\n", titleStyle.Render(m.title))

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return b.String()
}

func createTextInput(placeholder string) textinput.Model {
	t := textinput.New()
	t.Placeholder = placeholder
	t.CursorStyle = cursorStyle
	return t
}
