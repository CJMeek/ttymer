package tui

import (
	"fmt"
	"os"
	"time"

	"ttymer/task"

	tea "github.com/charmbracelet/bubbletea"
)

type sessionState int

const (
	taskView sessionState = iota
	taskEditView
	timerView
)

type MainModel struct {
	state        sessionState
	taskView     tea.Model
	taskEditView tea.Model
	timerView    tea.Model
	tasks        []*task.Task
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case EditTaskMsg:
		m.taskEditView = NewTaskEditModel(msg.Task)
		m.state = taskEditView
		return m, nil
	case DoneEditingMsg:
		m.state = taskView
		return m, nil
	case TimerViewMsg:
		total := time.Duration(0)
		for _, t := range m.tasks {
			total += t.Duration
		}
		m.timerView = NewTimerView(total)
		m.state = timerView
		return m, m.timerView.Init()
	case TimerBackMsg:
		m.state = taskView
		return m, nil
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch m.state {
	case taskView:
		newTaskView, newCmd := m.taskView.Update(msg)
		taskViewModel, ok := newTaskView.(TaskViewModel)
		if !ok {
			panic("Could not perform assertion on TaskViewModel")
		}
		m.taskView = taskViewModel
		cmd = newCmd
	case taskEditView:
		newTaskEditView, newCmd := m.taskEditView.Update(msg)
		taskEditViewModel, ok := newTaskEditView.(TaskEditModel)
		if !ok {
			panic("Could not perform assertion on TaskEditViewModel")
		}
		m.taskEditView = taskEditViewModel
		cmd = newCmd
	case timerView:
		newTimerView, newCmd := m.timerView.Update(msg)
		timerViewModel, ok := newTimerView.(TimerView)
		if !ok {
			panic("Could not perform assertion on TimerView model")
		}
		m.timerView = timerViewModel
		cmd = newCmd
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func InitTui() (tea.Model, tea.Cmd) {
	tasks := []*task.Task{
		{ID: 1, Name: "Example 1", TaskDescription: "Example 1", Duration: 1 * time.Minute},
		{ID: 2, Name: "Example 2", TaskDescription: "Example 2", Duration: 1 * time.Minute},
	}

	m := MainModel{
		state:    taskView,
		tasks:    tasks,
		taskView: NewTaskView(tasks),
	}
	return m, nil
}

func StartTea() {
	m, _ := InitTui()

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("broke", err)
		os.Exit(1)
	}
}

func (m MainModel) View() string {

	switch m.state {
	case taskView:
		return m.taskView.View()
	case taskEditView:
		return m.taskEditView.View()
	case timerView:
		return m.timerView.View()
	default:
		return m.taskView.View()
	}
}
