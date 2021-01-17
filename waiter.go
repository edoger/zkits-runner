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
	"sync"
)

// Waiter interface defines the waiter.
type Waiter interface {
	// Wait blocks the current coroutine and waits for the current waiter to be closed.
	// For waiters that have been closed, this method will not block.
	Wait()

	// Channel returns a read-only channel that can be used for select.
	// For waiters that have been closed, this method returns a closed channel.
	Channel() <-chan struct{}
}

// ReceiptableWaiter interface defines the receiptable waiter.
type ReceiptableWaiter interface {
	Waiter

	// Done reports that the current waiter has completed and is about to exit.
	Done()
}

// CloseableWaiter interface defines the closeable waiter.
type CloseableWaiter interface {
	Waiter

	// Waiter returns a pure waiter.
	Waiter() Waiter

	// Close closes the current waiter. This method is idempotent.
	Close()
}

// DuplexWaiter interface defines the duplex waiter.
type DuplexWaiter interface {
	ReceiptableWaiter

	// Waiter returns a pure receiptable waiter.
	Waiter() ReceiptableWaiter

	// Close closes the current waiter. This method is idempotent.
	Close()

	// WaitDone waits for the Done() method of the current waiter to be called.
	// This method returns only after the Done() method is called by the coroutine holding
	// the wait. This method is idempotent.
	// Essentially, this method is relative to: <-DoneChannel().
	WaitDone()

	// DoneChannel returns a read-only channel. When the Done() method of the current waiter
	// is called, this channel will be closed.
	DoneChannel() <-chan struct{}

	// CloseAndWaitDone closes the current waiter and waits for the Done method of the
	// current waiter to be called.
	CloseAndWaitDone()
}

// The built-in Waiter.
type channelWaiter struct {
	c chan struct{}
}

// Create and return a new built-in Waiter instance.
func newChannelWaiter() *channelWaiter {
	return &channelWaiter{c: make(chan struct{})}
}

// Wait blocks the current coroutine and waits for the current waiter to be closed.
// For waiters that have been closed, this method will not block.
// Essentially, this method is relative to: <-Channel().
func (w *channelWaiter) Wait() { <-w.Channel() }

// Channel returns a read-only channel that can be used for select.
// For waiters that have been closed, this method returns a closed channel.
func (w *channelWaiter) Channel() <-chan struct{} {
	return w.c
}

// NewCloseableWaiter creates and returns a new CloseableWaiter instance.
func NewCloseableWaiter() CloseableWaiter {
	return newCloseableWaiter()
}

// The built-in CloseableWaiter.
type closeableWaiter struct {
	*channelWaiter
	once sync.Once
}

// Create and return a new built-in CloseableWaiter instance.
func newCloseableWaiter() *closeableWaiter {
	return &closeableWaiter{channelWaiter: newChannelWaiter()}
}

// Waiter returns a pure waiter.
func (w *closeableWaiter) Waiter() Waiter { return w.channelWaiter }

// Close closes the current waiter. This method is idempotent.
func (w *closeableWaiter) Close() { w.once.Do(w.close) }

func (w *closeableWaiter) close() { close(w.c) }

// The built-in ReceiptableWaiter.
type receiptableWaiter struct {
	*channelWaiter
	d    chan struct{}
	once sync.Once
}

// Create and return a new built-in ReceiptableWaiter instance.
func newReceiptableWaiter() *receiptableWaiter {
	return &receiptableWaiter{channelWaiter: newChannelWaiter(), d: make(chan struct{})}
}

// Done reports that the current waiter has completed and is about to exit.
func (w *receiptableWaiter) Done() { w.once.Do(w.done) }

func (w *receiptableWaiter) done() { close(w.d) }

// NewDuplexWaiter creates and returns a new DuplexWaiter instance.
func NewDuplexWaiter() DuplexWaiter {
	return newDuplexWaiter()
}

// The built-in DuplexWaiter.
type duplexWaiter struct {
	*receiptableWaiter
	once sync.Once
}

// Create and return a new built-in DuplexWaiter instance.
func newDuplexWaiter() *duplexWaiter {
	return &duplexWaiter{
		receiptableWaiter: newReceiptableWaiter(),
	}
}

// Waiter returns a pure receiptable waiter.
func (w *duplexWaiter) Waiter() ReceiptableWaiter {
	return w.receiptableWaiter
}

// Close closes the current waiter. This method is idempotent.
func (w *duplexWaiter) Close() { w.once.Do(w.close) }

func (w *duplexWaiter) close() { close(w.c) }

// WaitDone waits for the Done() method of the current waiter to be called.
// This method returns only after the Done() method is called by the coroutine holding
// the wait. This method is idempotent.
// Essentially, this method is relative to: <-DoneChannel().
func (w *duplexWaiter) WaitDone() { <-w.DoneChannel() }

// DoneChannel returns a read-only channel. When the Done() method of the current waiter
// is called, this channel will be closed.
func (w *duplexWaiter) DoneChannel() <-chan struct{} {
	return w.d
}

// CloseAndWaitDone closes the current waiter and waits for the Done() method of the
// current waiter to be called.
func (w *duplexWaiter) CloseAndWaitDone() {
	w.Close()
	w.WaitDone()
}
