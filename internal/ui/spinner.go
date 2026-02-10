package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SpinnerModel is a Bubbletea model that displays a spinner with a message.
type SpinnerModel struct {
	spinner  spinner.Model
	message  string
	done     bool
	result   string
	err      error
	action   func() (string, error)
	quitting bool
}

type actionDoneMsg struct {
	result string
	err    error
}

// NewSpinner creates a new spinner model with a message and an action to run.
func NewSpinner(message string, action func() (string, error)) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(Purple)
	return SpinnerModel{
		spinner: s,
		message: message,
		action:  action,
	}
}

func (m SpinnerModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.runAction())
}

func (m SpinnerModel) runAction() tea.Cmd {
	return func() tea.Msg {
		result, err := m.action()
		return actionDoneMsg{result: result, err: err}
	}
}

func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	case actionDoneMsg:
		m.done = true
		m.result = msg.result
		m.err = msg.err
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m SpinnerModel) View() string {
	if m.quitting {
		return ""
	}
	if m.done {
		if m.err != nil {
			return ErrorBox.Render(Error.Render("Error: ") + m.err.Error()) + "\n"
		}
		return m.result
	}
	return fmt.Sprintf("%s %s\n", m.spinner.View(), m.message)
}

// RunSpinner runs a spinner with the given message while the action executes.
// Falls back to plain output when no TTY is available.
func RunSpinner(message string, action func() (string, error)) (string, error) {
	if !IsTTY() {
		fmt.Println(message)
		result, err := action()
		if err != nil {
			fmt.Println(Error.Render("Error: ") + err.Error())
			return "", err
		}
		return result, nil
	}

	model := NewSpinner(message, action)
	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}
	m := finalModel.(SpinnerModel)
	if m.err != nil {
		return "", m.err
	}
	return m.result, nil
}
