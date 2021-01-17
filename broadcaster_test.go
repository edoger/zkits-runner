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
	"testing"
)

func TestBroadcaster(t *testing.T) {
	b := NewBroadcaster()
	if b == nil {
		t.Fatal("NewBroadcaster(): nil")
	}

	ss := make([]string, 0)
	go func(w ReceiptableWaiter) {
		defer w.Done()
		w.Wait()
		ss = append(ss, "A")
	}(b.NewWaiter())

	go func(w ReceiptableWaiter) {
		defer w.Done()
		<-w.Channel()
		ss = append(ss, "B")
	}(b.NewWaiter())

	b.Broadcast()
	if got := strings.Join(ss, "-"); got != "B-A" {
		t.Fatalf("Broadcaster.Broadcast(): %s", got)
	}

	go func(w ReceiptableWaiter) {
		defer w.Done()
		w.Wait()
		ss = append(ss, "C")
	}(b.NewWaiter())

	go func(w ReceiptableWaiter) {
		defer w.Done()
		<-w.Channel()
		ss = append(ss, "D")
	}(b.NewWaiter())

	b.Close()
	if got := strings.Join(ss, "-"); got != "B-A-D-C" {
		t.Fatalf("Broadcaster.Close(): %s", got)
	}

	select {
	case <-b.NewWaiter().Channel():
		ss = append(ss, "F")
	default:
		ss = append(ss, "E")
	}
	if got := strings.Join(ss, "-"); got != "B-A-D-C-F" {
		t.Fatalf("Broadcaster.NewWaiter(): %s", got)
	}
}
