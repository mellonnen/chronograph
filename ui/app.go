package ui

import (
	"errors"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mellonnen/chronograph/models"
	"gorm.io/gorm"
)

type state int

const (
	loadWorspaces state = iota
	showWorkspaces
	showRepos
	showTasks
	showTaskInfo

	addWorkspace
	addRepo
	addTask

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
	return tea.NewProgram(m, tea.WithAltScreen())
}

func (m model) Init() tea.Cmd {
	m.state = showWaiting
	m.waitingText = "initializing the database..."
	return tea.Batch(m.spinner.Tick, initSqliteCmd(m.dbPath))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case dbMsg:
		m.db = msg.DB
		m.waitingText = "fetching workspaces"
		cmds = append(cmds, listWorkspacesCmd(m.db))

	case addWorkspaceMsg:
		res := m.db.Create(&msg.Workspace)
		if res.RowsAffected != 1 {
			return m, errorCmd(errors.New("create ineffective"))
		}
		m.workspaces = append(m.workspaces, msg.Workspace)

	case removeResourceMsg:
		switch m.state {
		case showWorkspaces:
			res := m.db.Unscoped().Delete(&m.workspaces[msg.index])
			if res.RowsAffected != 1 {
				return m, errorCmd(errors.New("remove ineffective"))
			}
			m.workspaces = append(m.workspaces[:msg.index], m.workspaces[msg.index+1:]...)
		}

	case listWorkspacesMsg:
		m.workspaces = msg.Workspaces
		m.list = newList(m.workspaces, Workspace)
		m.state = showWorkspaces

	case errorMsg:
		m.err = msg
		m.state = showError
	}

	switch m.state {
	case showWorkspaces, showRepos, showTasks:
		newList, cmd := m.list.update(msg)
		m.list = newList
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
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
