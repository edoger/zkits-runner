// Copyright 2020 The ZKits Project Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package runner

import (
	"errors"
	"sync"
)

// ErrExited returns when running a task in an exited Runner.
var ErrExited = errors.New("runner: exited")

// Runner defines the task runner.
// The task runner is used to manage the operation and shutdown of multiple
// independent subtasks of an application.
type Runner interface {
	// Run method executes the given task instance synchronously.
	// If the runner has exited, the ErrExited error will be returned.
	Run(Task) error

	// MustRun method executes the given task instance synchronously.
	// If the task execution returns a non nil error, panic immediately.
	MustRun(Task) Runner

	// Wait method blocks the current coroutine until the runner exits.
	// When the exit signal is received or the exit method is called,
	// the blocking state of the method is released.
	Wait() error

	// WaitBy method blocks the current coroutine until the runner exits.
	// When a given channel is closed or the exit method is called,
	// the blocking state of the method is released.
	WaitBy(<-chan struct{}) error

	// Exit method exits the current runner.
	Exit() error

	// Exited method determines whether the current runner has exited.
	Exited() bool
}

// New creates and returns a new instance of the Runner.
func New() Runner {
	return &runner{chanExit: make(chan struct{})}
}

// The runner type is an implementation of the built-in Runner.
type runner struct {
	mutex    sync.Mutex
	tasks    []Task
	chanExit chan struct{}
	onceExit sync.Once
}

// Run method executes the given task instance synchronously.
// If the runner has exited, the ErrExited error will be returned.
func (r *runner) Run(t Task) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.Exited() {
		return ErrExited
	}

	if err := SafeCall(t.Execute); err != nil {
		return err
	}
	r.tasks = append(r.tasks, t)
	return nil
}

// MustRun method executes the given task instance synchronously.
// If the task execution returns a non nil error, panic immediately.
func (r *runner) MustRun(t Task) Runner {
	MustCall(func() error { return r.Run(t) })
	return r
}

// Wait method blocks the current coroutine until the runner exits.
// When the exit signal is received or the exit method is called,
// the blocking state of the method is released.
func (r *runner) Wait() error {
	return r.WaitBy(GetSystemExitChan())
}

// WaitBy method blocks the current coroutine until the runner exits.
// When a given channel is closed or the exit method is called,
// the blocking state of the method is released.
func (r *runner) WaitBy(c <-chan struct{}) error {
	select {
	case <-c:
		return r.Exit()
	case <-r.chanExit:
		// In this case, because the Exit method is called, do nothing!
		return nil
	}
}

// Exit method exits the current runner.
func (r *runner) Exit() error {
	r.mutex.Lock()
	defer func() {
		// Make sure to unblock the Wait method.
		r.onceExit.Do(func() { close(r.chanExit) })
		r.mutex.Unlock()
	}()
	// In this case, we don't care about the state of the runner, just
	// make sure that all tasks in the current runner are shut down.
	if len(r.tasks) == 0 {
		return nil
	}

	err := new(Errors)
	for i := len(r.tasks) - 1; i >= 0; i-- {
		err.Add(SafeCall(r.tasks[i].Shutdown))
	}
	r.tasks = r.tasks[:0]

	if err.Len() > 1 {
		return err
	}
	return err.First()
}

// Exited method determines whether the current runner has exited.
func (r *runner) Exited() bool {
	select {
	case <-r.chanExit:
		return true
	default:
		return false
	}
}
