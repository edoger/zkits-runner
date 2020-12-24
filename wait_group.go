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

// WaitGroup waits for a collection of goroutines to finish.
// We only provide WaitGroup.Go and WaitGroup.MultiGo on the basis of sync.WaitGroup.
// Example:
//   var wg WaitGroup
//   wg.Go(func() { /* do something */ })
//   wg.Wait()
// or
//   wg.MultiGo(5, func() { /* do something */ })
//   wg.Wait()
type WaitGroup struct {
	wg sync.WaitGroup
}

// Go uses a goroutines to run the f function.
func (w *WaitGroup) Go(f func()) {
	w.wg.Add(1)
	go w.do(f)
}

// MultiGo uses n goroutines to run the f function. Panic if n <= 0.
func (w *WaitGroup) MultiGo(n int, f func()) {
	if n <= 0 {
		panic("WaitGroup.MultiGo(): n must be a positive integer")
	}

	w.wg.Add(n)
	for i := 0; i < n; i++ {
		go w.do(f)
	}
}

func (w *WaitGroup) do(f func()) {
	defer w.wg.Done()
	f()
}

// Wait blocks waiting for all goroutines to exit.
func (w *WaitGroup) Wait() {
	w.wg.Wait()
}
