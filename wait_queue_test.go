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
	"strings"
	"sync"
	"testing"
)

func TestWaitQueue(t *testing.T) {
	wq := NewWaitQueue()
	if wq == nil {
		t.Fatal("NewWaitQueue(): nil")
	}

	wg1 := new(sync.WaitGroup)
	wg2 := new(sync.WaitGroup)
	wg1.Add(5)
	wg2.Add(5)
	ss := make([]string, 5)

	go func(w Waiter) {
		wg1.Done()
		w.Wait()
		ss[0] = "test1"
		wg2.Done()
	}(wq.NewWaiter())
	go func(w Waiter) {
		wg1.Done()
		w.Wait()
		ss[1] = "test2"
		wg2.Done()
	}(wq.NewWaiter())
	go func(w Waiter) {
		wg1.Done()
		w.Wait()
		ss[2] = "test3"
		wg2.Done()
	}(wq.NewWaiter())
	go func(w Waiter) {
		wg1.Done()
		w.Wait()
		ss[3] = "test4"
		wg2.Done()
	}(wq.NewWaiter())
	go func(w Waiter) {
		wg1.Done()
		w.Wait()
		ss[4] = "test5"
		wg2.Done()
	}(wq.NewWaiter())

	wg1.Wait()

	if n := wq.Len(); n != 5 {
		t.Fatalf("WaitQueue.Len(): %d", n)
	}

	if n := wq.Release(2); n != 2 {
		t.Fatalf("WaitQueue.Release(): %d", n)
	} else {
		if l := wq.Len(); l != 3 {
			t.Fatalf("WaitQueue.Release(): %d", l)
		}
	}

	if n := wq.ReleaseAll(); n != 3 {
		t.Fatalf("WaitQueue.ReleaseAll(): %d", n)
	} else {
		if l := wq.Len(); l != 0 {
			t.Fatalf("WaitQueue.ReleaseAll(): %d", l)
		}
	}

	wg2.Wait()
	if s := strings.Join(ss, ""); s != "test1test2test3test4test5" {
		t.Fatalf("WaitQueue: %s", s)
	}
}
