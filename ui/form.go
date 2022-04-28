package ui

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

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
	inputs     []inputModel
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

	m.inputs = make([]inputModel, 0)
	m.inputs = append(m.inputs, newInput("Name", func(s string) bool { return len(s) > 0 }))
	m.inputs = append(m.inputs, newInput("Description"))

	switch r {
	case Repo:
		pathInput := newInput("Path to Repo")
		cwd, _ := os.Getwd()
		path, _ := git.RepoPathFromPath(cwd)
		pathInput.Input.SetValue(path)
		m.inputs = append(m.inputs, pathInput)
	case Task:
		validate := func(s string) bool {
			_, err := time.ParseDuration(s)
			if err != nil {
				return false
			}
			return true
		}
		m.inputs = append(m.inputs, newInput("Estimated time", validate))
	}

	m.inputs[0].Input.Focus()
	m.inputs[0].Input.PromptStyle = focusedStyle
	m.inputs[0].Input.TextStyle = focusedStyle
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
				// check for invalid fields.
				valid := true
				for _, input := range m.inputs {
					if input.valid != nil && !*input.valid {
						valid = false
					}
				}
				if !valid {
					break
				}
				switch m.resource {
				case Workspace:
					workspace := models.Workspace{
						Name:        m.inputs[0].Input.Value(),
						Description: sql.NullString{String: m.inputs[1].Input.Value(), Valid: true},
					}

					return m, addWorkspaceCmd(workspace)
				case Repo:
					path := m.inputs[2].Input.Value()
					remote, err := git.RemoteFromPath(path)
					if err != nil {
						return m, errorCmd(fmt.Errorf("getting remote: %v", err))
					}
					repo := models.Repo{
						Name:        m.inputs[0].Input.Value(),
						Description: sql.NullString{String: m.inputs[1].Input.Value(), Valid: true},
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
					cmds[i] = m.inputs[i].Input.Focus()
					m.inputs[i].Input.PromptStyle = focusedStyle
					m.inputs[i].Input.TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Input.Blur()
				m.inputs[i].Input.PromptStyle = noStyle
				m.inputs[i].Input.TextStyle = noStyle
			}
			return m, tea.Batch(cmds...)
		}
	}
	// Update cursor blinking.
	var cmds = make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m formModel) view() string {
	var b strings.Builder

	fmt.Fprintf(&b, "%s\n\n", titleStyle.Render(m.title))

	for i := range m.inputs {
		b.WriteString(m.inputs[i].view())
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

type validationFunc func(string) bool
type inputModel struct {
	Input    textinput.Model
	valid    *bool
	validate validationFunc
}

func newInput(placeholder string, validate ...validationFunc) inputModel {
	m := inputModel{}
	m.Input = textinput.New()
	m.Input.Placeholder = placeholder
	m.Input.CursorStyle = cursorStyle
	if len(validate) < 1 {
		m.validate = func(s string) bool { return true }
	} else {
		m.validate = validate[0]
	}
	return m
}

func (m inputModel) update(msg tea.Msg) (inputModel, tea.Cmd) {
	m.valid = boolPtr(m.validate(m.Input.Value()))
	if len(m.Input.Value()) == 0 && *m.valid {
		m.valid = nil
	}
	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

func (m inputModel) view() string {
	var valid rune
	switch {
	case m.valid == nil:
		valid = '🟡'
	case *m.valid == true:
		valid = '🟢'
	case *m.valid == false:
		valid = '🔴'
	}
	return fmt.Sprintf("%c %s", valid, m.Input.View())
}

func boolPtr(b bool) *bool {
	return &b
}
