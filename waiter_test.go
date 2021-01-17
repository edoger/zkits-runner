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
	"testing"
	"time"
)

func TestCloseableWaiter(t *testing.T) {
	waiter := NewCloseableWaiter()
	if waiter == nil {
		t.Fatal("NewCloseableWaiter(): nil")
	}

	wg := new(sync.WaitGroup)
	wg.Add(4)

	var m, n, p, q int

	go func(w Waiter) {
		defer wg.Done()
		w.Wait()
		m = 1
	}(waiter.Waiter())

	go func() {
		defer wg.Done()
		<-waiter.Channel()
		n = 1
	}()

	go func(w Waiter) {
		defer wg.Done()
		<-w.Channel()
		p = 1
	}(waiter.Waiter())

	go func() {
		defer wg.Done()
		waiter.Wait()
		q = 1
	}()

	time.Sleep(time.Millisecond * 100)

	waiter.Close()
	wg.Wait()

	if m != 1 || n != 1 || p != 1 || q != 1 {
		t.Fatalf("CloseableWaiter: %d %d %d %d", m, n, p, q)
	}
}

func TestDuplexWaiter(t *testing.T) {
	waiter := NewDuplexWaiter()
	if waiter == nil {
		t.Fatal("NewDuplexWaiter(): nil")
	}

	wg := new(sync.WaitGroup)
	wg.Add(4)

	var m, n, p, q int

	go func(w ReceiptableWaiter) {
		defer wg.Done()
		w.Wait()
		m = 1
		w.Done()
	}(waiter.Waiter())
	go func(w ReceiptableWaiter) {
		defer wg.Done()
		<-w.Channel()
		n = 1
		w.Done()
	}(waiter.Waiter())

	go func() {
		defer wg.Done()
		waiter.Wait()
		p = 1
		waiter.Done()
	}()
	go func() {
		defer wg.Done()
		<-waiter.Channel()
		q = 1
		waiter.Done()
	}()

	time.Sleep(time.Millisecond * 100)

	waiter.CloseAndWaitDone()
	wg.Wait()

	if m != 1 || n != 1 || p != 1 || q != 1 {
		t.Fatalf("NewDuplexWaiter: %d %d %d %d", m, n, p, q)
	}
}

// func TestBroadcaster(t *testing.T) {
// 	b := NewBroadcaster()
// 	if b == nil {
// 		t.Fatal("NewBroadcaster(): nil")
// 	}
//
// 	var ss []string
//
// 	go func(w Waiter) {
// 		defer w.Done()
// 		w.Wait()
// 		ss = append(ss, "A")
// 	}(b.NewWaiter())
//
// 	go func(w Waiter) {
// 		defer w.Done()
// 		<-w.Channel()
// 		ss = append(ss, "B")
// 	}(b.NewWaiter())
//
// 	b.Broadcast()
// 	if got := strings.Join(ss, "-"); got != "B-A" {
// 		t.Fatalf("Broadcaster.Broadcast(): %s", got)
// 	}
//
// 	go func(w Waiter) {
// 		defer w.Done()
// 		w.Wait()
// 		ss = append(ss, "C")
// 	}(b.NewWaiter())
//
// 	go func(w Waiter) {
// 		defer w.Done()
// 		<-w.Channel()
// 		ss = append(ss, "D")
// 	}(b.NewWaiter())
//
// 	b.Close()
// 	if got := strings.Join(ss, "-"); got != "B-A-D-C" {
// 		t.Fatalf("Broadcaster.Close(): %s", got)
// 	}
//
// 	select {
// 	case <-b.NewWaiter().Channel():
// 		ss = append(ss, "F")
// 	default:
// 		ss = append(ss, "E")
// 	}
// 	if got := strings.Join(ss, "-"); got != "B-A-D-C-F" {
// 		t.Fatalf("Broadcaster.NewWaiter(): %s", got)
// 	}
// }
