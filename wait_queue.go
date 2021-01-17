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

// WaitQueue interface defines the wait queue.
type WaitQueue interface {
	// NewWaiter creates a waiter and adds it to the wait queue.
	NewWaiter() Waiter

	// Len returns the number of waiters in the current queue.
	Len() int

	// Release releases up to the top n waiters in the queue.
	// This method returns the number of released waiters, the range is [0, n].
	// The release sequence is the same as the enqueue sequence.
	Release(int) int

	// ReleaseAll releases all the waiters in the queue.
	// This method returns the number of released waiters.
	// The release sequence is the same as the enqueue sequence.
	ReleaseAll() int
}

// NewWaitQueue creates and returns a new WaitQueue instance.
func NewWaitQueue() WaitQueue {
	return new(waitQueue)
}

// The built-in WaitQueue.
type waitQueue struct {
	mutex sync.Mutex
	queue []CloseableWaiter
}

// NewWaiter creates a waiter and adds it to the wait queue.
func (wq *waitQueue) NewWaiter() Waiter {
	wq.mutex.Lock()
	defer wq.mutex.Unlock()

	w := NewCloseableWaiter()
	wq.queue = append(wq.queue, w)
	return w.Waiter()
}

// Len returns the number of waiters in the current queue.
func (wq *waitQueue) Len() (n int) {
	wq.mutex.Lock()
	n = len(wq.queue)
	wq.mutex.Unlock()
	return
}

// Release releases up to the top n waiters in the queue.
// This method returns the number of released waiters, the range is [0, n].
// The release sequence is the same as the enqueue sequence.
func (wq *waitQueue) Release(n int) int {
	wq.mutex.Lock()
	defer wq.mutex.Unlock()

	if max := len(wq.queue); max > 0 && n > 0 {
		var queue []CloseableWaiter
		for i, j := max-1, max-n; i >= 0; i-- {
			if i < j {
				queue = wq.queue[:i+1]
				break
			}
			wq.queue[i].Close()
		}
		wq.queue = queue
		return max - len(queue)
	}
	return 0
}

// ReleaseAll releases all the waiters in the queue.
// This method returns the number of released waiters.
// The release sequence is the same as the enqueue sequence.
func (wq *waitQueue) ReleaseAll() (n int) {
	wq.mutex.Lock()
	defer wq.mutex.Unlock()

	if n = len(wq.queue); n > 0 {
		for i := n - 1; i >= 0; i-- {
			wq.queue[i].Close()
		}
		wq.queue = nil
	}
	return
}
