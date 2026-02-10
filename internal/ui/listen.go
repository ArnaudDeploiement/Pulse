package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"pulse/internal/transport"
)

// ListenModel is the Bubbletea model for the listen view.
type ListenModel struct {
	spinner   spinner.Model
	groupName string
	storeDir  string
	events    <-chan transport.ReceiveEvent
	received  []fileEntry
	errors    []string
	quitting  bool
	startTime time.Time
}

type fileEntry struct {
	name string
	size int64
	from string
	at   time.Time
}

type receiveEventMsg transport.ReceiveEvent

// NewListenModel creates a listener UI.
func NewListenModel(groupName, storeDir string, events <-chan transport.ReceiveEvent) ListenModel {
	s := spinner.New()
	s.Spinner = spinner.Pulse
	s.Style = lipgloss.NewStyle().Foreground(Cyan)
	return ListenModel{
		spinner:   s,
		groupName: groupName,
		storeDir:  storeDir,
		events:    events,
		received:  make([]fileEntry, 0),
		errors:    make([]string, 0),
		startTime: time.Now(),
	}
}

func (m ListenModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.waitForEvent())
}

func (m ListenModel) waitForEvent() tea.Cmd {
	return func() tea.Msg {
		ev, ok := <-m.events
		if !ok {
			return tea.Quit()
		}
		return receiveEventMsg(ev)
	}
}

func (m ListenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	case receiveEventMsg:
		ev := transport.ReceiveEvent(msg)
		if ev.Err != nil {
			m.errors = append(m.errors, ev.Err.Error())
		} else {
			m.received = append(m.received, fileEntry{
				name: ev.Filename,
				size: ev.Size,
				from: ev.From,
				at:   time.Now(),
			})
		}
		return m, m.waitForEvent()
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m ListenModel) View() string {
	if m.quitting {
		return Muted.Render("Listener stopped.") + "\n"
	}

	header := lipgloss.JoinHorizontal(lipgloss.Center,
		m.spinner.View(),
		" ",
		Subtitle.Render(fmt.Sprintf("Listening on group %q", m.groupName)),
	)

	uptime := time.Since(m.startTime).Round(time.Second)
	info := Muted.Render(fmt.Sprintf("  Store: %s  |  Uptime: %s  |  Files: %d",
		m.storeDir, uptime, len(m.received)))

	s := "\n" + header + "\n" + info + "\n\n"

	if len(m.received) > 0 {
		s += Subtitle.Render("  Received files:") + "\n"
		// Show last 10 entries
		start := 0
		if len(m.received) > 10 {
			start = len(m.received) - 10
		}
		for _, f := range m.received[start:] {
			short := f.from
			if len(short) > 16 {
				short = short[:8] + "..." + short[len(short)-8:]
			}
			s += fmt.Sprintf("  %s %s  %s  %s\n",
				Success.Render("[OK]"),
				Highlight.Render(f.name),
				Muted.Render(formatSize(f.size)),
				Muted.Render("from "+short),
			)
		}
		s += "\n"
	}

	if len(m.errors) > 0 {
		s += Warning.Render("  Errors:") + "\n"
		start := 0
		if len(m.errors) > 5 {
			start = len(m.errors) - 5
		}
		for _, e := range m.errors[start:] {
			s += fmt.Sprintf("  %s %s\n", Error.Render("[ERR]"), Muted.Render(e))
		}
		s += "\n"
	}

	s += Muted.Render("  Press q or Ctrl+C to stop") + "\n"

	return s
}

func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// RunListener runs the listener UI.
// Falls back to plain output when no TTY is available.
func RunListener(groupName, storeDir string, events <-chan transport.ReceiveEvent) error {
	if !IsTTY() {
		fmt.Printf("Listening on group %q -> %s\n", groupName, storeDir)
		for ev := range events {
			if ev.Err != nil {
				fmt.Printf("[ERR] %s\n", ev.Err)
			} else {
				fmt.Printf("[OK] %s (%s) from %s\n", ev.Filename, formatSize(ev.Size), ev.From)
			}
		}
		return nil
	}

	model := NewListenModel(groupName, storeDir, events)
	p := tea.NewProgram(model)
	_, err := p.Run()
	return err
}
