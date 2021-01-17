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
	"github.com/edoger/zkits-runner/internal"
)

// This is the global default empty wait.
// Since the behavior of the empty wait is deterministic, we share an empty wait instance.
var globalEmptyReceiptableWaiter = new(emptyReceiptableWaiter)

// EmptyReceiptableWaiter returns an empty receiptable waiter instance.
// For the returned empty waiter, the Wait method does not block, the Channel method
// always returns a closed channel, and the Done method is an empty function.
func EmptyReceiptableWaiter() ReceiptableWaiter { return globalEmptyReceiptableWaiter }

// The emptyWaiter type defines an empty waiter.
type emptyReceiptableWaiter struct{}

// Wait implements the Waiter interface, but do nothing.
func (*emptyReceiptableWaiter) Wait() { /* Do nothing */ }

// Channel returns a closed read-only channel.
func (*emptyReceiptableWaiter) Channel() <-chan struct{} {
	return internal.ClosedChan()
}

// Done implements the Waiter interface, but do nothing.
func (*emptyReceiptableWaiter) Done() { /* Do nothing */ }
