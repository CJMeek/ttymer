package task

import (
	"fmt"
	"time"
)

type Task struct {
	ID              int
	Name            string
	TaskDescription string
	Duration        time.Duration
}

func (t Task) Title() string { return t.Name }

func (t Task) Description() string {
	if t.Duration > 0 {
		return fmt.Sprintf("%s • %s", t.TaskDescription, t.Duration)
	}
	return t.TaskDescription
}

func (t Task) FilterValue() string { return t.Name }
