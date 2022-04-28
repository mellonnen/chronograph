package ui

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mellonnen/chronograph/models"
	"gorm.io/gorm"
)

var appStyle = lipgloss.NewStyle().Padding(1, 2)

type state int

const (
	loadWorspaces state = iota
	showWorkspaces
	showRepos
	showTasks
	showTaskInfo

	showCreateWorkspace
	showCreateRepo
	showCreateTask

	showWaiting
	showError
)

type model struct {
	state state

	list listModel
	form formModel

	workspaces       []models.Workspace
	currentWorkspace *models.Workspace
	currentRepo      *models.Repo
	currentTask      *models.Task

	dbPath string
	db     *gorm.DB

	height int
	width  int

	err         error
	waitingText string
}

func New(dbPath string) *tea.Program {
	m := model{dbPath: dbPath}
	return tea.NewProgram(m, tea.WithAltScreen())
}

func (m model) Init() tea.Cmd {
	return tea.Batch(initSqliteCmd(m.dbPath))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case dbMsg:
		m.db = msg.DB
		m.waitingText = "fetching workspaces"
		cmds = append(cmds, listWorkspacesCmd(m.db))

	case createResourceMsg:
		switch m.state {
		case showWorkspaces:
			m.state = showCreateWorkspace
			m.form = newForm(Workspace)
		case showRepos:
			m.state = showCreateRepo
			m.form = newForm(Repo)
		}
		cmds = append(cmds, m.form.init())

	case removeResourceMsg:
		switch m.state {
		case showWorkspaces:
			res := m.db.Unscoped().Delete(&m.workspaces[msg.index])
			if res.RowsAffected != 1 {
				return m, errorCmd(errors.New("remove ineffective"))
			}
			m.workspaces = append(m.workspaces[:msg.index], m.workspaces[msg.index+1:]...)
		}
	case addWorkspaceMsg:
		// add workspace to database.
		res := m.db.Create(&msg.Workspace)
		if res.RowsAffected != 1 {
			return m, errorCmd(errors.New("create ineffective"))
		}
		m.workspaces = append(m.workspaces, msg.Workspace)
		cmds = append(cmds, addResourceCmd(msg.Workspace))
		m.state = showWorkspaces

	case addRepoMsg:
		err := m.db.Model(m.currentWorkspace).Association("Repos").Append(&msg.Repo)
		if err == nil {
			return m, errorCmd(err)
		}
		m.db.Preload("Repos").Find(m.currentWorkspace)
		cmds = append(cmds, addResourceCmd(msg.Repo))
		m.state = showRepos

	case chooseResourceMsg:
		switch m.state {
		case showWorkspaces:
			m.currentWorkspace = &m.workspaces[msg.index]
			m.db.Preload("Repos").Find(m.currentWorkspace)
			m.list = newList(m.currentWorkspace.Repos, Repo, m.height, m.width)
			m.state = showRepos
		}

	case listWorkspacesMsg:
		m.workspaces = msg.Workspaces
		m.list = newList(m.workspaces, Workspace, m.height, m.width)
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
	case showCreateWorkspace, showCreateRepo, showCreateTask:
		newForm, cmd := m.form.update(msg)
		m.form = newForm
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	switch m.state {
	case showError:
		return m.errorView()
	case showWorkspaces, showRepos, showTasks:
		return appStyle.Render(m.list.view())
	case showCreateWorkspace, showCreateRepo, showCreateTask:
		return appStyle.Render(m.form.view())
	default:
		return ""
	}
}
