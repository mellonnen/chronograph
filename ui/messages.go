package ui

import (
	"github.com/mellonnen/chronograph/models"
	"gorm.io/gorm"
)

type errorMsg error

type dbMsg struct {
	DB *gorm.DB
}

type listWorkspacesMsg struct {
	Workspaces []models.Workspace
}

type addWorkspaceMsg struct {
	Workspace models.Workspace
}

type addRepoMsg struct {
	Repo models.Repo
}

type createResourceMsg struct{}
type addResourceMsg struct {
	Resource models.Listable
}

type removeResourceMsg struct {
	index int
}

type chooseResourceMsg struct {
	index int
}
