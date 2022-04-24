package ui

import (
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

	err
	waiting
)

type model struct {
	state uiState

	spinner spinner.Model

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
	s.Spinner = spinner.MiniDot
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
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case dbMsg:
		m.db = msg.DB
		m.state = root
		m.waitingText = ""
		return m, nil

	case errorMsg:
		m.err = msg
		m.state = err
		return m, nil
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	switch m.state {
	case waiting:
		return m.waitingView()
	case err:
		return m.errorView()
	default:
		return m.rootView()
	}
}
