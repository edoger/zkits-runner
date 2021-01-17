// Copyright 2021 The ZKits Project Authors.
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

// Broadcaster interface defines the broadcaster.
type Broadcaster interface {
	// NewWaiter creates and returns a new Waiter instance.
	// It is worth noting that the order of closing the wait is opposite to
	// that of creation, and the closing process is linear.
	// The Waiter returned by this method is one-time, and once it is closed,
	// it will always be closed. If the broadcaster is closed, then this method
	// will always return an empty waiter.
	NewWaiter() ReceiptableWaiter

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
	waiters []DuplexWaiter
	closed  bool
}

// NewWaiter creates and returns a new Waiter instance.
// It is worth noting that the order of closing the wait is opposite to
// that of creation, and the closing process is linear.
// The Waiter returned by this method is one-time, and once it is closed,
// it will always be closed. If the broadcaster is closed, then this method
// will always return an empty waiter.
func (b *broadcaster) NewWaiter() ReceiptableWaiter {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.closed {
		return EmptyReceiptableWaiter()
	}
	w := NewDuplexWaiter()
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
			b.waiters[i].CloseAndWaitDone()
		}
		b.waiters = nil
	}
}
