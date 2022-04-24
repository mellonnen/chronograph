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
