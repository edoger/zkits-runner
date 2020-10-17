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
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestNewTaskFromFunc(t *testing.T) {
	var tasks []Task

	tasks = []Task{
		NewTaskFromFunc(nil),
		NewTaskFromFunc(func() error { return nil }),
		NewTaskFromFunc(func() error { return nil }, func() error { return nil }),
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

func TestNewTaskFromHTTPServer(t *testing.T) {
	task := NewTaskFromHTTPServer(new(http.Server), nil)
	if task == nil {
		t.Fatal("NewTaskFromHTTPServer(): nil task")
	}
}

func TestNewTaskFromHTTPServerPanic(t *testing.T) {
	var v interface{}
	func() {
		defer func() { v = recover() }()
		NewTaskFromHTTPServer(nil, nil)
	}()

	if v == nil {
		t.Fatal("NewTaskFromHTTPServer(): no panic")
	}
}

func TestHTTPServerTask(t *testing.T) {
	addr := "127.0.0.1:38281"
	task := NewTaskFromHTTPServer(&http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = io.WriteString(w, "ok")
		}),
	}, nil)

	r := New()
	if err := r.Run(task); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond * 20)

	if res, err := http.Get("http://" + addr); err != nil {
		t.Fatal(err)
	} else {
		data, err := ioutil.ReadAll(res.Body)
		_ = res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		if s := string(data); s != "ok" {
			t.Fatalf("The http server response: %s", s)
		}
	}

	if err := r.Exit(); err != nil {
		t.Fatal(err)
	}
	if err := task.Shutdown(); err != nil {
		t.Fatal(err)
	}
	if err := task.Execute(); err != http.ErrServerClosed {
		t.Fatalf("The http server task not closed: %s", err)
	}
}
