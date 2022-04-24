package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type Workspace struct {
	gorm.Model
	Name        string
	Description sql.NullString

	Repos []Repo
}

type Repo struct {
	gorm.Model
	WorkspaceID uint

	Name        string
	Description sql.NullString
	Remote      string
	Path        string

	Tasks []Task
}

type Task struct {
	gorm.Model
	RepoID      uint
	Name        string
	Description sql.NullString

	StartedAt   sql.NullTime
	CompletedAt sql.NullTime

	StartSHA []byte
	EndSHA   []byte
}
