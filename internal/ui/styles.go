package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	Gold   = lipgloss.Color("#FFD700")
	Onyx   = lipgloss.Color("#1A1A1A")
	Bronze = lipgloss.Color("#CD7F32")
	Gray   = lipgloss.Color("#666666")
	Cyan   = lipgloss.Color("#00FFFF")
	Red    = lipgloss.Color("#FF4500")

	// Base Styles
	AppStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Foreground(Gold)

	TitleStyle = lipgloss.NewStyle().
			Foreground(Onyx).
			Background(Gold).
			Padding(0, 1).
			Bold(true).
			MarginBottom(1)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Bronze).
			Padding(0, 1)

	ShadowStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#333333"))

	LabelStyle = lipgloss.NewStyle().
			Foreground(Gray).
			Italic(true)

	ValueStyle = lipgloss.NewStyle().
			Foreground(Gold).
			Bold(true)

	HighlightStyle = lipgloss.NewStyle().
			Foreground(Cyan)

	AlertStyle = lipgloss.NewStyle().
			Foreground(Red).
			Bold(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(Gray).
			MarginTop(1)
)

func DrawBox(title string, content string, width int) string {
	box := BoxStyle.Width(width).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			TitleStyle.Render(" " + title + " "),
			content,
		),
	)
	return box
}
