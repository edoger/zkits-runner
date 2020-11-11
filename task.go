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
	"context"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

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

// NewTaskFromHTTPServer creates and returns a Task from a given http.Server.
// If the given http.Server is nil, panic immediately.
func NewTaskFromHTTPServer(server *http.Server, options *HTTPServerOptions) Task {
	if server == nil {
		panic("NewTaskFromHTTPServer(): nil http server")
	}
	if options == nil {
		options = new(HTTPServerOptions)
	}
	return &httpServerTask{server: server, options: options}
}

// HTTPServerOptions defines options when the http.Server is run as a Task.
type HTTPServerOptions struct {
	// ErrorHandler is used to handle http.Server errors.
	// Normally, it is executed at most once during the running period of this http.Server.
	ErrorHandler func(error)

	// TLSCertFile and TLSKeyFile are used to enable TLS and are only valid if they are not empty.
	TLSCertFile, TLSKeyFile string
}

// The httpServerTask type is used to wrap the http.Server as a task.
type httpServerTask struct {
	mutex   sync.Mutex
	status  int64
	err     error
	server  *http.Server
	options *HTTPServerOptions
}

// Execute starts the current http server task.
// This method is idempotent. If an error occurs in the http server, this method
// will always return the last error.
func (t *httpServerTask) Execute() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if atomic.CompareAndSwapInt64(&t.status, 0, 1) {
		wg := new(sync.WaitGroup)
		wg.Add(1)
		go t.execute(wg)
		wg.Wait()
		// Sleep for 5 milliseconds to ensure that the http server
		// is started as much as possible.
		time.Sleep(time.Millisecond * 5)
	}
	return t.err
}

// The execute method will block waiting for the http server to exit.
func (t *httpServerTask) execute(wg *sync.WaitGroup) {
	// Report immediately that the goroutine has started running.
	go wg.Done()

	if t.options.TLSCertFile != "" && t.options.TLSKeyFile != "" {
		t.err = t.server.ListenAndServeTLS(t.options.TLSCertFile, t.options.TLSKeyFile)
	} else {
		t.err = t.server.ListenAndServe()
	}

	if t.err != nil {
		// Since the http.ErrServerClosed error is returned when the http server
		// exits normally, we will not report it here.
		if t.err != http.ErrServerClosed && t.options.ErrorHandler != nil {
			t.options.ErrorHandler(t.err)
		}
	}
}

// Shutdown shuts down the http server.
// This method is idempotent.
func (t *httpServerTask) Shutdown() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if atomic.CompareAndSwapInt64(&t.status, 1, 2) {
		return t.server.Shutdown(context.Background())
	}
	return nil
}
