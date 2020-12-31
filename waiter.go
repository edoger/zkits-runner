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

	"github.com/edoger/zkits-runner/internal"
)

// The Waiter interface defines the waiter.
type Waiter interface {
	// Wait blocks the current coroutine and waits for the current waiter to be closed.
	// For waiters that have been closed, this method will not block.
	// Essentially, this method is relative to: <-Channel().
	Wait()

	// Channel returns a read-only channel that can be used for select.
	// For waiters that have been closed, this method returns a closed channel.
	Channel() <-chan struct{}

	// Done reports that the current coroutine will exit.
	// In an asynchronous coroutine, be sure to call this method before exiting.
	Done()
}

// This is the global default empty wait.
// Since the behavior of the empty wait is deterministic, we share an empty wait instance.
var defaultEmptyWaiter = newEmptyWaiter()

// GetEmptyWaiter returns an empty waiter instance.
// For the returned empty waiter, the Wait method does not block, the Channel method
// always returns a closed channel, and the Done method is an empty function.
func GetEmptyWaiter() Waiter {
	return defaultEmptyWaiter
}

// The newEmptyWaiter function creates and returns an empty waiter instance.
func newEmptyWaiter() Waiter {
	return &emptyWaiter{internal.ClosedChan()}
}

// The emptyWaiter type defines an empty waiter.
type emptyWaiter struct {
	c <-chan struct{}
}

// Wait does nothing.
func (w *emptyWaiter) Wait() { /* Do nothing */ }

// Channel returns a closed channel.
func (w *emptyWaiter) Channel() <-chan struct{} { return w.c }

// Done does nothing.
func (w *emptyWaiter) Done() { /* Do nothing */ }

// The newChannelWaiter function returns a new instance of the built-in waiter.
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
// For waiters that have been closed, this method will not block.
// Essentially, this method is relative to: <-Channel().
func (w *channelWaiter) Wait() {
	<-w.messageChan
}

// Channel returns a read-only channel that can be used for select.
// For waiters that have been closed, this method returns a closed channel.
func (w *channelWaiter) Channel() <-chan struct{} {
	return w.messageChan
}

// Done reports that the current coroutine will exit.
// In an asynchronous coroutine, be sure to call this method before exiting.
func (w *channelWaiter) Done() {
	w.receiptOnce.Do(w.closeReceiptChan)
}

// Close the receipt channel.
func (w *channelWaiter) closeReceiptChan() {
	close(w.receiptChan)
}

// Close the current waiter and waits for the Done method to be called.
func (w *channelWaiter) close() {
	w.messageOnce.Do(w.closeMessageChan)
	<-w.receiptChan
}

// Close the message channel.
func (w *channelWaiter) closeMessageChan() {
	close(w.messageChan)
}

// The CloseableWaiter interface defines the closeable waiter.
type CloseableWaiter interface {
	Waiter

	// Waiter returns a pure waiter.
	// The behavior of the wait returned by this method is consistent with that of
	// itself, except that the shutdown control function is hidden, which helps to
	// ensure that the wait is not accidentally closed in the child coroutine.
	Waiter() Waiter

	// Close closes the current waiter and waits for the Done method to be called.
	// This method returns only after the Waiter.Done() method is called by the
	// coroutine holding the wait. This method is idempotent.
	Close()
}

// NewCloseableWaiter creates and returns a new CloseableWaiter instance.
func NewCloseableWaiter() CloseableWaiter {
	return &channelCloseableWaiter{newChannelWaiter()}
}

// The built-in implementation of the CloseableWaiter interface.
type channelCloseableWaiter struct {
	*channelWaiter
}

// Waiter returns a pure waiter.
// The behavior of the wait returned by this method is consistent with that of
// itself, except that the shutdown control function is hidden, which helps to
// ensure that the wait is not accidentally closed in the child coroutine.
func (w *channelCloseableWaiter) Waiter() Waiter {
	return w.channelWaiter
}

// Close closes the current waiter and waits for the Done method to be called.
// This method returns only after the Waiter.Done method is called by the
// coroutine holding the wait. This method is idempotent.
func (w *channelCloseableWaiter) Close() {
	w.close()
}

// The Broadcaster interface defines the broadcaster.
type Broadcaster interface {
	// NewWaiter creates and returns a new Waiter instance.
	// It is worth noting that the order of closing the wait is opposite to
	// that of creation, and the closing process is linear.
	// The Waiter returned by this method is one-time, and once it is closed,
	// it will always be closed. If the broadcaster is closed, then this method
	// will always return an empty waiter.
	NewWaiter() Waiter

	// Broadcast sends a close signal to all the waiters that have been created
	// and waits for all the waiters to call the Waiter.Done method.
	// After this method is called, the broadcaster will return to its initial state.
	Broadcast()

	// Close closes the current broadcaster.
	// The behavior of this method is consistent with the Broadcast method, the only
	// difference is that after this method returns, the NewWaiter method will always
	// return an empty waiter instance.
	Close()
}

// NewBroadcaster creates and returns a new broadcaster instance.
func NewBroadcaster() Broadcaster {
	return &broadcaster{}
}

// The built-in implementation of the Broadcaster interface.
type broadcaster struct {
	mutex   sync.Mutex
	waiters []CloseableWaiter
	closed  bool
}

// NewWaiter creates and returns a new Waiter instance.
// It is worth noting that the order of closing the wait is opposite to
// that of creation, and the closing process is linear.
// The Waiter returned by this method is one-time, and once it is closed,
// it will always be closed. If the broadcaster is closed, then this method
// will always return an empty waiter.
func (b *broadcaster) NewWaiter() Waiter {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.closed {
		return GetEmptyWaiter()
	}
	w := NewCloseableWaiter()
	b.waiters = append(b.waiters, w)
	return w.Waiter()
}

// Broadcast sends a close signal to all the waiters that have been created
// and waits for all the waiters to call the Waiter.Done method.
// After this method is called, the broadcaster will return to its initial state.
func (b *broadcaster) Broadcast() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.close()
}

// Close closes the current broadcaster.
// The behavior of this method is consistent with the Broadcast method, the only
// difference is that after this method returns, the NewWaiter method will always
// return an empty waiter instance.
func (b *broadcaster) Close() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.closed = true
	b.close()
}

// Close all the waiters in the current broadcaster in reverse order.
func (b *broadcaster) close() {
	if n := len(b.waiters); n > 0 {
		for i := len(b.waiters) - 1; i >= 0; i-- {
			b.waiters[i].Close()
		}
		b.waiters = nil
	}
}
