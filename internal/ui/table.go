package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Table renders a simple styled table.
type Table struct {
	Headers []string
	Rows    [][]string
}

// Render produces the styled table string.
func (t Table) Render() string {
	if len(t.Headers) == 0 {
		return ""
	}

	// Calculate column widths
	widths := make([]int, len(t.Headers))
	for i, h := range t.Headers {
		widths[i] = len(h)
	}
	for _, row := range t.Rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Add padding
	for i := range widths {
		widths[i] += 3
	}

	var sb strings.Builder

	// Header
	headerRow := ""
	for i, h := range t.Headers {
		headerRow += TableHeader.Width(widths[i]).Render(h)
	}
	sb.WriteString("  " + headerRow + "\n")

	// Separator
	sep := "  "
	for _, w := range widths {
		sep += lipgloss.NewStyle().Foreground(Dim).Render(strings.Repeat("─", w))
	}
	sb.WriteString(sep + "\n")

	// Rows
	for _, row := range t.Rows {
		rowStr := ""
		for i, cell := range row {
			if i < len(widths) {
				style := TableCell
				if i > 0 {
					style = TableCellDim
				}
				rowStr += style.Width(widths[i]).Render(cell)
			}
		}
		sb.WriteString("  " + rowStr + "\n")
	}

	if len(t.Rows) == 0 {
		sb.WriteString("  " + Muted.Render("(empty)") + "\n")
	}

	return sb.String()
}

// KeyValue renders a key-value pair nicely.
func KeyValue(key, value string) string {
	return fmt.Sprintf("  %s %s", Subtitle.Render(key+":"), value)
}
