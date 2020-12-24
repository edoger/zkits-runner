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
	"sync/atomic"
	"testing"
)

func TestWaitGroup(t *testing.T) {
	var wg WaitGroup
	var n int64

	wg.Go(func() { atomic.AddInt64(&n, 1) })
	wg.MultiGo(5, func() { atomic.AddInt64(&n, 1) })
	wg.Wait()

	if n != 6 {
		t.Fatalf("WaitGroup: %d", n)
	}
}

func TestWaitGroup_Panic(t *testing.T) {
	defer func() {
		if v := recover(); v == nil {
			t.Fatal("WaitGroup: no panic")
		}
	}()
	new(WaitGroup).MultiGo(0, func() {})
}
