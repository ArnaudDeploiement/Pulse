package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Color palette
	Purple    = lipgloss.Color("#7C3AED")
	Cyan      = lipgloss.Color("#06B6D4")
	Green     = lipgloss.Color("#10B981")
	Red       = lipgloss.Color("#EF4444")
	Yellow    = lipgloss.Color("#F59E0B")
	Dim       = lipgloss.Color("#6B7280")
	White     = lipgloss.Color("#F9FAFB")
	DarkGray  = lipgloss.Color("#1F2937")
	LightGray = lipgloss.Color("#9CA3AF")

	// Base styles
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(Purple).
		MarginBottom(1)

	Subtitle = lipgloss.NewStyle().
			Foreground(Cyan).
			Bold(true)

	Success = lipgloss.NewStyle().
		Foreground(Green).
		Bold(true)

	Error = lipgloss.NewStyle().
		Foreground(Red).
		Bold(true)

	Warning = lipgloss.NewStyle().
		Foreground(Yellow)

	Muted = lipgloss.NewStyle().
		Foreground(Dim)

	Highlight = lipgloss.NewStyle().
			Foreground(Cyan)

	// Box styles
	InfoBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Purple).
		Padding(1, 2).
		MarginTop(1).
		MarginBottom(1)

	SuccessBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Green).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	ErrorBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Red).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	// Table styles
	TableHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(Purple).
			PaddingRight(2)

	TableCell = lipgloss.NewStyle().
			Foreground(White).
			PaddingRight(2)

	TableCellDim = lipgloss.NewStyle().
			Foreground(LightGray).
			PaddingRight(2)

	// Badges
	BadgeActive = lipgloss.NewStyle().
			Background(Green).
			Foreground(lipgloss.Color("#000000")).
			Padding(0, 1).
			Bold(true)

	BadgeInactive = lipgloss.NewStyle().
			Background(Dim).
			Foreground(White).
			Padding(0, 1)

	// Key hint style (for bottom bar)
	KeyStyle = lipgloss.NewStyle().
			Foreground(Cyan).
			Bold(true)

	DescStyle = lipgloss.NewStyle().
			Foreground(Dim)
)
