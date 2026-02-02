package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time

type TimerView struct {
	total    time.Duration
	elapsed  time.Duration
	running  bool
	lastTick time.Time
	bar      progress.Model
}

func NewTimerView(total time.Duration) TimerView {
	bar := progress.New(
		progress.WithSolidFill("█"),
		progress.WithWidth(40),
	)
	return TimerView{
		total:   total,
		running: true,
		bar:     bar,
	}
}

func (m TimerView) Init() tea.Cmd { return tick() }

func tick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (m TimerView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.bar.Width = msg.Width - 4
		if m.bar.Width < 10 {
			m.bar.Width = 10
		}
		return m, nil
	case tickMsg:
		t := time.Time(msg)
		if m.running {
			if !m.lastTick.IsZero() {
				m.elapsed += t.Sub(m.lastTick)
			}
			m.lastTick = t
		}
		percent := 0.0
		if m.total > 0 {
			percent = float64(m.elapsed) / float64(m.total)
			if percent > 1 {
				percent = 1
			}
		}
		cmd := m.bar.SetPercent(percent)
		return m, tea.Batch(tick(), cmd)
	case progress.FrameMsg:
		bar, cmd := m.bar.Update(msg)
		m.bar = bar.(progress.Model)
		return m, cmd
	case tea.KeyMsg:
		switch msg.String() {
		case "p":
			m.running = !m.running
			if !m.running {
				m.lastTick = time.Time{}
			}
		}
	}
	return m, nil
}

func (m TimerView) View() string {
	if m.total <= 0 {
		return "Timer\n\nNo duration set.\n"
	}

	remaining := m.total - m.elapsed
	if remaining < 0 {
		remaining = 0
	}

	percent := 0.0
	if m.total > 0 {
		percent = float64(m.elapsed) / float64(m.total)
		if percent > 1 {
			percent = 1
		}
	}

	return "Timer\n\n" + m.bar.View() + "\n\nElapsed: " + m.elapsed.Round(time.Second).String() + "\nRemaining: " + remaining.Round(time.Second).String() + "\nPercent: " + fmt.Sprintf("%.1f%%", percent*100) + "\n"

}
