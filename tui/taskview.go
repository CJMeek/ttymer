package tui

import (
	"ttymer/task"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type EditTaskMsg struct {
	Task *task.Task
}

type TimerViewMsg struct{}

type keymap struct {
	Create key.Binding
	Edit   key.Binding
	Delete key.Binding
	Timer  key.Binding
}

var keyMap = keymap{
	Create: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "create"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Timer: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "go to timer"),
	),
}

type TaskViewModel struct {
	taskList list.Model
}

func NewTaskView(tasks []*task.Task) tea.Model {
	items := make([]list.Item, 0, len(tasks))
	for _, t := range tasks {
		items = append(items, t)
	}

	m := TaskViewModel{taskList: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.taskList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			keyMap.Create,
			keyMap.Delete,
			keyMap.Edit,
			keyMap.Timer,
		}
	}

	return m
}

func (m TaskViewModel) Init() tea.Cmd {
	return nil
}

func (m TaskViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.taskList.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "e":
			if selected, ok := m.taskList.SelectedItem().(*task.Task); ok {
				return m, func() tea.Msg { return EditTaskMsg{Task: selected} }
			}
		case "t":
			return m, func() tea.Msg { return TimerViewMsg{} }
		}
	}

	var cmd tea.Cmd
	m.taskList, cmd = m.taskList.Update(msg)
	return m, cmd
}

func (m TaskViewModel) View() string {

	return docStyle.Render(m.taskList.View())

}
