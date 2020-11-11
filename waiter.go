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

// The Waiter interface defines the waiter.
type Waiter interface {
	// Wait blocks the current coroutine and waits for the current waiter to be closed.
	Wait()

	// Channel returns a read-only channel that can be used for select.
	Channel() <-chan struct{}

	// Done reports that the current coroutine will exit.
	Done()
}

// The CloseableWaiter interface defines the closeable waiter.
type CloseableWaiter interface {
	Waiter

	// Waiter returns a pure waiter.
	Waiter() Waiter

	// Close closes the current waiter and waits for the Done method to be called.
	Close()
}

// NewCloseableWaiter creates and returns a new CloseableWaiter instance.
func NewCloseableWaiter() CloseableWaiter {
	return &channelCloseableWaiter{newChannelWaiter()}
}

// The channelCloseableWaiter type is a built-in implementation of the CloseableWaiter interface.
type channelCloseableWaiter struct {
	*channelWaiter
}

// Waiter returns a pure waiter.
func (w *channelCloseableWaiter) Waiter() Waiter {
	return w.channelWaiter
}

// Close closes the current waiter and waits for the Done method to be called.
func (w *channelCloseableWaiter) Close() {
	w.close()
}

// The newChannelWaiter function creates and returns a new instance of the built-in waiter.
func newChannelWaiter() *channelWaiter {
	return &channelWaiter{
		messageChan: make(chan struct{}),
		receiptChan: make(chan struct{}),
	}
}

// The channelWaiter type is a built-in implementation of the Waiter interface.
type channelWaiter struct {
	messageChan, receiptChan chan struct{}
	messageOnce, receiptOnce sync.Once
}

// Wait blocks the current coroutine and waits for the current waiter to be closed.
func (w *channelWaiter) Wait() {
	<-w.messageChan
}

// Channel returns a read-only channel that can be used for select.
func (w *channelWaiter) Channel() <-chan struct{} {
	return w.messageChan
}

// Done reports that the current coroutine will exit.
func (w *channelWaiter) Done() {
	w.receiptOnce.Do(w.closeReceiptChan)
}

// This method closes the receipt channel.
func (w *channelWaiter) closeReceiptChan() {
	close(w.receiptChan)
}

// Close the current waiter and waits for the Done method to be called.
func (w *channelWaiter) close() {
	w.messageOnce.Do(w.closeMessageChan)
	<-w.receiptChan
}

// This method closes the message channel.
func (w *channelWaiter) closeMessageChan() {
	close(w.messageChan)
}
