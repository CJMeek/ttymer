package tui

import (
	"strings"
	"time"

	"ttymer/task"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type DoneEditingMsg struct{}

type TaskEditModel struct {
	task          *task.Task
	nameInput     textinput.Model
	descInput     textinput.Model
	durationInput textinput.Model
	focusIndex    int
	errMsg        string
}

func NewTaskEditModel(t *task.Task) tea.Model {
	nameInput := textinput.New()
	nameInput.Placeholder = "Task name"
	nameInput.SetValue(t.Name)
	nameInput.Focus()

	descInput := textinput.New()
	descInput.Placeholder = "Description"
	descInput.SetValue(t.TaskDescription)

	durationInput := textinput.New()
	durationInput.Placeholder = "Duration (e.g. 25m, 1h30m)"
	if t.Duration > 0 {
		durationInput.SetValue(t.Duration.String())
	}

	return TaskEditModel{
		task:          t,
		nameInput:     nameInput,
		descInput:     descInput,
		durationInput: durationInput,
	}
}

func (m TaskEditModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m TaskEditModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return DoneEditingMsg{} }
		case "tab", "shift+tab", "up", "down":
			if msg.String() == "up" || msg.String() == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > 2 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = 2
			}

			if m.focusIndex == 0 {
				m.nameInput.Focus()
				m.descInput.Blur()
				m.durationInput.Blur()
			} else {
				m.nameInput.Blur()
				if m.focusIndex == 1 {
					m.descInput.Focus()
					m.durationInput.Blur()
				} else {
					m.descInput.Blur()
					m.durationInput.Focus()
				}
			}
		case "enter":
			if m.focusIndex == 2 {
				durationValue := strings.TrimSpace(m.durationInput.Value())
				if durationValue != "" {
					parsed, err := time.ParseDuration(durationValue)
					if err != nil {
						m.errMsg = "Invalid duration. Try 25m, 1h30m, 2h."
						return m, nil
					}
					m.task.Duration = parsed
				} else {
					m.task.Duration = 0
				}

				m.task.Name = m.nameInput.Value()
				m.task.TaskDescription = m.descInput.Value()
				m.errMsg = ""
				return m, func() tea.Msg { return DoneEditingMsg{} }
			}

			m.focusIndex++
			if m.focusIndex > 2 {
				m.focusIndex = 2
			}
			if m.focusIndex == 1 {
				m.nameInput.Blur()
				m.descInput.Focus()
				m.durationInput.Blur()
			} else if m.focusIndex == 2 {
				m.nameInput.Blur()
				m.descInput.Blur()
				m.durationInput.Focus()
			}
		}
	}

	var cmd tea.Cmd
	m.nameInput, cmd = m.nameInput.Update(msg)
	m.descInput, _ = m.descInput.Update(msg)
	m.durationInput, _ = m.durationInput.Update(msg)
	return m, cmd
}

func (m TaskEditModel) View() string {
	var b strings.Builder
	b.WriteString("Edit task\n\n")
	b.WriteString("Name:\n")
	b.WriteString(m.nameInput.View())
	b.WriteString("\n\nDescription:\n")
	b.WriteString(m.descInput.View())
	b.WriteString("\n\nDuration:\n")
	b.WriteString(m.durationInput.View())
	if m.errMsg != "" {
		b.WriteString("\n\n")
		b.WriteString(m.errMsg)
	}
	b.WriteString("\n\nEnter to save, Esc to cancel")
	return b.String()
}
