package ui

import (
	"strings"
	"time"

	"github.com/mellonnen/chronograph/models"
)

const timeFmt = "2006-01-02 15:04:05"

type overviewModel struct {
	task models.Task
}

func newOverwiew(task models.Task) overviewModel {
	return overviewModel{
		task: task,
	}
}

func (m overviewModel) view() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render(m.task.Name))
	b.WriteString("\n\n")

	if m.task.Description.Valid {
		b.WriteString(primaryStyle.Render("Description: "))
		b.WriteString(secondaryStyle.Render(m.task.Description.String))
		b.WriteString("\n\n")
	}

	b.WriteString(primaryStyle.Render("Status: "))
	var status string
	if m.task.CompletedAt.Valid {
		status = "Complete"
	} else {
		status = "Incomplete"
	}
	b.WriteString(secondaryStyle.Render(status))
	b.WriteString("\n\n")

	b.WriteString(primaryStyle.Render("Estimated time: "))
	b.WriteString(secondaryStyle.Render(shortDur(m.task.ExpectedDuration)))
	b.WriteString("\n\n")

	b.WriteString(primaryStyle.Render("Created: "))
	b.WriteString(secondaryStyle.Render(m.task.CreatedAt.Format(timeFmt)))
	b.WriteString("\n\n")

	b.WriteString(primaryStyle.Render("Updated: "))
	b.WriteString(secondaryStyle.Render(m.task.UpdatedAt.Format(timeFmt)))
	b.WriteString("\n\n")

	return b.String()
}

func shortDur(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
}
