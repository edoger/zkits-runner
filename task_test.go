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
	"testing"
)

func TestNewTaskFromFunc(t *testing.T) {
	var tasks []Task

	tasks = []Task{
		NewTaskFromFunc(nil),
		NewTaskFromFunc(func() error {
			return nil
		}),
		NewTaskFromFunc(func() error {
			return nil
		}, func() error {
			return nil
		}),
	}
	for _, task := range tasks {
		if task == nil {
			t.Fatal("NewTaskFromFunc(): nil task")
		}
		if err := task.Execute(); err != nil {
			t.Fatalf("Task.Execute(): %s", err)
		}
		if err := task.Shutdown(); err != nil {
			t.Fatalf("Task.Shutdown(): %s", err)
		}
	}
}

func TestNewTaskFromFuncPanic(t *testing.T) {
	var v interface{}
	func() {
		defer func() { v = recover() }()
		NewTaskFromFunc(nil, nil, nil)
	}()

	if v == nil {
		t.Fatal("NewTaskFromFunc(nil, nil, nil): no panic")
	}
}
