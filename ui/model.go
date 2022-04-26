package ui

import (
	"database/sql"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mellonnen/chronograph/models"
	"gorm.io/gorm"
)

type uiState int

const (
	root uiState = iota
	inWorkspace
	inRepo
	inTask

	waiting
	showError
)

type model struct {
	state uiState

	spinner spinner.Model

	list             list.Model
	listKeys         listKeyMap
	workspaces       []models.Workspace
	currentWorkspace *models.Workspace
	currentRepo      *models.Repo
	currentTask      *models.Task

	dbPath string
	db     *gorm.DB

	err         error
	waitingText string
}

func New(dbPath string) *tea.Program {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	m := model{spinner: s, dbPath: dbPath}
	return tea.NewProgram(m)
}

func (m model) Init() tea.Cmd {
	m.state = waiting
	m.waitingText = "initializing the database..."
	return tea.Batch(m.spinner.Tick, initSqliteCmd(m.dbPath))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			return m, nil
		}
		switch {
		case msg.String() == "q", msg.String() == "ctrl-c":
			return m, tea.Quit
		case key.Matches(msg, m.listKeys.create):
			workspace := models.Workspace{
				Name:        "Misc",
				Description: sql.NullString{String: "General workspace", Valid: true},
			}
			return m, tea.Batch(addWorkspaceCmd(m.db, workspace))
		default:
			return m, nil
		}
	case dbMsg:
		m.db = msg.DB
		m.waitingText = "fetching workspaces"
		return m, tea.Batch(listWorkspacesCmd(m.db))

	case listWorkspacesMsg:
		m.workspaces = msg.Workspaces
		m.list = newList(m.workspaces)
		m.list.Title = "Workspaces"
		m.listKeys = newListKeyMap("workspace")
		m.list.AdditionalFullHelpKeys = func() []key.Binding {
			return []key.Binding{
				m.listKeys.create,
				m.listKeys.choose,
				m.listKeys.remove,
				m.listKeys.toggleHelp,
			}
		}
		m.state = root
		return m, nil

	case errorMsg:
		m.err = msg
		m.state = showError
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	switch m.state {
	case waiting:
		return m.waitingView()
	case showError:
		return m.errorView()
	default:
		return m.rootView()
	}
}
