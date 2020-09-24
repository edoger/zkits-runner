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

// Task interface defines the task units that the runner can run.
type Task interface {
	// Execute method is the entry point for the task to run.
	// When the task is run by the runner, this method is executed first.
	Execute() error

	// Shutdown method is the method to exit the task.
	// When the exit is run, this method will be called.
	Shutdown() error
}

// NewTaskFromFunc creates a runnable task from a given function.
func NewTaskFromFunc(execute func() error, shutdown ...func() error) Task {
	switch len(shutdown) {
	case 0:
		return &funcTask{execute: execute}
	case 1:
		return &funcTask{execute: execute, shutdown: shutdown[0]}
	default:
		panic("NewTaskFromFunc(): too many shutdown function.")
	}
}

// The funcTask type is used to wrap a given function into a runnable task.
type funcTask struct {
	execute, shutdown func() error
}

// Execute method executes the given execute function.
// If the given function is nil, ignored.
func (t *funcTask) Execute() error {
	if t.execute == nil {
		return nil
	}
	return t.execute()
}

// Shutdown method executes the given shutdown function.
// If the given function is nil, ignored.
func (t *funcTask) Shutdown() error {
	if t.shutdown == nil {
		return nil
	}
	return t.shutdown()
}
