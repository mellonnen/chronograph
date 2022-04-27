package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mellonnen/chronograph/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initSqliteCmd(dbPath string) tea.Cmd {
	return func() tea.Msg {
		db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			return errorMsg(fmt.Errorf("initializing sqlite database: %w", err))
		}
		db.AutoMigrate(&models.Workspace{}, &models.Repo{}, &models.Task{})
		return dbMsg{DB: db}
	}
}

func listWorkspacesCmd(db *gorm.DB) tea.Cmd {
	return func() tea.Msg {
		var workspaces []models.Workspace
		db.Find(&workspaces)
		return listWorkspacesMsg{Workspaces: workspaces}
	}
}
func addWorkspaceCmd(workspace models.Workspace) tea.Cmd {
	return func() tea.Msg {
		return addWorkspaceMsg{Workspace: workspace}
	}
}

func createResourceCmd() tea.Cmd {
	return func() tea.Msg {
		return createResourceMsg{}
	}
}

func addResourceCmd(resource models.Listable) tea.Cmd {
	return func() tea.Msg {
		return addResourceMsg{
			Resource: resource,
		}
	}
}

func removeResourceCmd(index int) tea.Cmd {
	return func() tea.Msg {
		return removeResourceMsg{index: index}
	}
}

func errorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return errorMsg(err)
	}
}
