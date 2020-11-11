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
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	chanExit = make(chan struct{})
	onceExit = new(sync.Once)
	onceWait = new(sync.Once)
)

// GetSystemExitChan function returns the system exit signal channel.
func GetSystemExitChan() <-chan struct{} {
	onceWait.Do(doWaitSystemExitSignal)
	return chanExit
}

// Start a coroutine to run the waitSystemExitSignal function.
func doWaitSystemExitSignal() {
	go waitSystemExitSignal()
}

// Wait the exit signal of the operating system.
func waitSystemExitSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(c)

	select {
	case <-c:
		onceExit.Do(func() { close(chanExit) })
	case <-chanExit:
		return
	}
}

// WaitSystemExit function will block the current coroutine until
// the exit signal of the operating system is captured.
func WaitSystemExit() {
	<-GetSystemExitChan()
}
