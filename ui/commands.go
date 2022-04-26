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

func addWorkspaceCmd(db *gorm.DB, workspace models.Workspace) tea.Cmd {
	return func() tea.Msg {
		res := db.Create(&workspace)
		if res.RowsAffected != 1 {
			return errorMsg(errors.New("create ineffective"))
		}
		return addWorkspaceMsg{Workspace: workspace}
	}
}
