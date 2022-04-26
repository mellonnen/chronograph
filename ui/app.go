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

type state int

const (
	showWorkspaces state = iota
	showRepos
	showTasks

	showTaskInfo

	showWaiting
	showError
)

type model struct {
	state state

	spinner spinner.Model
	list    listModel

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
	m.state = showWaiting
	m.waitingText = "initializing the database..."
	return tea.Batch(m.spinner.Tick, initSqliteCmd(m.dbPath))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.list.filterState() == list.Filtering {
			return m, nil
		}
		switch {
		case msg.String() == "q", msg.String() == "ctrl-c":
			return m, tea.Quit

		case key.Matches(msg, m.list.Keys.create):
			workspace := models.Workspace{
				Name:        "Misc",
				Description: sql.NullString{String: "General workspace", Valid: true},
			}
			return m, addWorkspaceCmd(m.db, workspace)

		case key.Matches(msg, m.list.Keys.remove):
			index := m.list.index()
			return m, deleteWorkspaceCmd(m.db, m.workspaces[index], index)
		default:
			return m, nil
		}

	case dbMsg:
		m.db = msg.DB
		m.waitingText = "fetching workspaces"
		return m, listWorkspacesCmd(m.db)

	case addWorkspaceMsg:
		m.workspaces = append(m.workspaces, msg.Workspace)
		m.list = newList(m.workspaces, Workspace)
		return m, nil

	case deleteWorkspaceMsg:
		m.workspaces = append(m.workspaces[:msg.index], m.workspaces[msg.index+1:]...)
		m.list = newList(m.workspaces, Workspace)
		return m, nil

	case listWorkspacesMsg:
		m.workspaces = msg.Workspaces
		m.list = newList(m.workspaces, Workspace)
		m.state = showWorkspaces
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
	case showWaiting:
		return m.waitingView()
	case showError:
		return m.errorView()
	default:
		return m.list.view()
	}
}
