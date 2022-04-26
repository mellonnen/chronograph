package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Listable interface {
	GetName() string
	GetDescription() string
	GetCreatedAt() string
}

type Workspace struct {
	gorm.Model
	Name        string
	Description sql.NullString

	Repos []Repo
}

func (w Workspace) GetName() string { return w.Name }
func (w Workspace) GetDescription() string {
	if w.Description.Valid {
		return w.Description.String
	}
	return "No description available"
}
func (w Workspace) GetCreatedAt() string { return w.CreatedAt.Format("January 2, 2006, 15:30") }

type Repo struct {
	gorm.Model
	WorkspaceID uint

	Name        string
	Description sql.NullString
	Remote      string
	Path        string

	Tasks []Task
}

func (r Repo) GetName() string { return r.Name }
func (r Repo) GetDescription() string {
	if r.Description.Valid {
		return r.Description.String
	}
	return "No description available"
}
func (r Repo) GetCreatedAt() string { return r.CreatedAt.Format("January 2, 2006, 15:30") }

type Task struct {
	gorm.Model
	RepoID      uint
	Name        string
	Description sql.NullString

	StartedAt        sql.NullTime
	CompletedAt      sql.NullTime
	ExpectedDuration time.Duration

	StartSHA []byte
	EndSHA   []byte
}

func (t Task) GetName() string { return t.Name }
func (t Task) GetDescription() string {
	if t.Description.Valid {
		return t.Description.String
	}
	return "No description available"
}
func (t Task) GetCreatedAt() string { return t.CreatedAt.Format("January 2, 2006, 15:30") }
