package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

// ProgressModel shows a progress bar for multi-peer file sending.
type ProgressModel struct {
	progress progress.Model
	total    int
	done     int
	results  []PeerResult
	channel  <-chan PeerResult
	finished bool
}

// PeerResult holds the result of sending to a single peer.
type PeerResult struct {
	PeerID string
	Ok     bool
	Err    error
}

type peerResultMsg PeerResult
type allDoneMsg struct{}

// NewProgress creates a progress bar model.
func NewProgress(total int, ch <-chan PeerResult) ProgressModel {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(50),
	)
	return ProgressModel{
		progress: p,
		total:    total,
		channel:  ch,
		results:  make([]PeerResult, 0, total),
	}
}

func (m ProgressModel) Init() tea.Cmd {
	return m.waitForResult()
}

func (m ProgressModel) waitForResult() tea.Cmd {
	return func() tea.Msg {
		result, ok := <-m.channel
		if !ok {
			return allDoneMsg{}
		}
		return peerResultMsg(result)
	}
}

func (m ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case peerResultMsg:
		m.done++
		m.results = append(m.results, PeerResult(msg))
		if m.done >= m.total {
			m.finished = true
			return m, tea.Quit
		}
		return m, m.waitForResult()
	case allDoneMsg:
		m.finished = true
		return m, tea.Quit
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}
	return m, nil
}

func (m ProgressModel) View() string {
	if m.total == 0 {
		return Muted.Render("No peers to send to.") + "\n"
	}

	pct := float64(m.done) / float64(m.total)
	bar := m.progress.ViewAs(pct)

	s := fmt.Sprintf("\n  Sending to %d peer(s)\n\n  %s  %d/%d\n\n",
		m.total, bar, m.done, m.total)

	for _, r := range m.results {
		short := r.PeerID
		if len(short) > 16 {
			short = short[:8] + "..." + short[len(short)-8:]
		}
		if r.Ok {
			s += fmt.Sprintf("  %s %s\n", Success.Render("[OK]"), short)
		} else {
			errMsg := "unknown error"
			if r.Err != nil {
				errMsg = r.Err.Error()
			}
			s += fmt.Sprintf("  %s %s  %s\n", Error.Render("[FAIL]"), short, Muted.Render(errMsg))
		}
	}

	if m.finished {
		okCount := 0
		for _, r := range m.results {
			if r.Ok {
				okCount++
			}
		}
		s += "\n" + SuccessBox.Render(fmt.Sprintf("Transfer complete: %d/%d successful", okCount, m.total)) + "\n"
	}

	return s
}

// RunProgress runs the progress bar UI until all results are received.
// Falls back to plain output when no TTY is available.
func RunProgress(total int, ch <-chan PeerResult) ([]PeerResult, error) {
	if !IsTTY() {
		var results []PeerResult
		for r := range ch {
			status := Success.Render("[OK]")
			if !r.Ok {
				status = Error.Render("[FAIL]")
			}
			short := r.PeerID
			if len(short) > 16 {
				short = short[:8] + "..." + short[len(short)-8:]
			}
			fmt.Printf("  %s %s\n", status, short)
			results = append(results, r)
		}
		return results, nil
	}

	model := NewProgress(total, ch)
	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}
	m := finalModel.(ProgressModel)
	return m.results, nil
}
