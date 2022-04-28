package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var titleStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("62")).
	Foreground(lipgloss.Color("230")).
	Padding(0, 1)

var primaryStyle = lipgloss.NewStyle().
	Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
	Padding(0, 0, 0, 2)

var secondaryStyle = primaryStyle.Copy().
	Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

var primarySelectedStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder(), false, false, false, true).
	BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
	Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
	Padding(0, 0, 0, 1)

var secondarySelectedStyle = primarySelectedStyle.Copy().
	Foreground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"})

var primaryDimmedStyle = lipgloss.NewStyle().
	Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
	Padding(0, 0, 0, 2)

var secondaryDimmedStyle = primaryDimmedStyle.Copy().
	Foreground(lipgloss.AdaptiveColor{Light: "#C2B8C2", Dark: "#4D4D4D"})
