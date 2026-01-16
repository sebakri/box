package cmd

import "github.com/charmbracelet/lipgloss"

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	warnStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	titleStyle   = lipgloss.NewStyle().Bold(true).MarginBottom(1)
	typeStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	binStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("246"))
	dimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	toolStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
)
