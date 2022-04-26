package ui

import (
	"fmt"
)

func (m model) waitingView() string {
	return fmt.Sprintf("\n\n%s %s", m.spinner.View(), m.waitingText)
}

func (m model) errorView() string {
	return fmt.Sprintf("An error occurred, please file an issue at https://github.com/mellonnen/chronograph \n\n Error Trace:\n%s", m.err.Error())
}
