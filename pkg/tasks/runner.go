package tasks

import (
	"context"
	"sync"
)

// TaskStatus is given to Runner callbacks to indicate the status of the tasks
// it is running.
type TasksStatus struct {
	Name    string
	Running bool
}

// Task is implemented by concurrent tasks a programme is expected to execute.
type Task interface {
	// TaskName returns the name of the Task for monitoring purposes.
	TaskName() string

	// RunTask is called by a Runner to execute it's job, and is expected to
	// run until either the context is canceled or an error is encountered,
	// afterwhich the runner will shut down.
	RunTask(context.Context) error
}

// Runner manages the execution of concurrent Tasks.
type Runner struct {
	// TaskStarting is called before each Task is started by the Runner.
	TaskStarting func(*TasksStatus)

	// TaskStopped is called after a Task stops gracefully.
	TaskStopped func(*TasksStatus)

	// TaskFailed is called after a Task returns an unexpected error.
	TaskFailed func(*TasksStatus, error)

	tasks []Task
}

// Add attaches a Task to a Runner, to be run when Run is called.
func (r *Runner) Add(t Task) {
	r.tasks = append(r.tasks, t)
}

// Run starts all the Tasks that have been Added to the Runner, until the
// context is canceled or a Task returns an unexpected error.
func (r *Runner) Run(ctx context.Context) error {
	wg := new(sync.WaitGroup)

	for _, task := range r.tasks {
		wg.Add(1)

		go func(ctx context.Context, wg *sync.WaitGroup, task Task) {
			defer wg.Done()

			if r.TaskStarting != nil {
				r.TaskStarting(&TasksStatus{Name: task.TaskName(), Running: true})
			}

			err := task.RunTask(ctx)
			if err != nil {
				if r.TaskFailed != nil {
					r.TaskFailed(&TasksStatus{Name: task.TaskName(), Running: false}, err)
				}
			} else {
				if r.TaskStopped != nil {
					r.TaskStopped(&TasksStatus{Name: task.TaskName(), Running: false})
				}
			}
		}(ctx, wg, task)
	}

	// TODO(jc): shutdown other tasks when one fails.

	wg.Wait()
	return nil
}
