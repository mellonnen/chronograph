package ui

import (
	"errors"
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

func removeResourceCmd(index int) tea.Cmd {
	return func() tea.Msg {
		return removeResourceMsg{index: index}
	}
}

func removeWorkspaceFromDBCmd(db *gorm.DB, workspace models.Workspace, index int) tea.Cmd {
	return func() tea.Msg {
		res := db.Unscoped().Delete(&workspace)
		if res.RowsAffected != 1 {
			return errorMsg(errors.New("delete ineffective"))
		}
		return removedResourceMsg{index: index}
	}
}

func errorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return errorMsg(err)
	}
}
