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
	"errors"
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	if r := New(); r == nil {
		t.Fatal("New(): return nil")
	}
}

func TestRunner_Run(t *testing.T) {
	r := New()

	if err := r.Run(NewTaskFromFunc(nil)); err != nil {
		t.Fatalf("Runner.Run(): %s", err)
	}
	if err := r.Run(NewTaskFromFunc(func() error { return errors.New("test1") })); err == nil {
		t.Fatal("Runner.Run(): nil error")
	}

	n := 0
	task1 := NewTaskFromFunc(func() error {
		n = 1
		return nil
	})
	if err := r.Run(task1); err != nil {
		t.Fatalf("Runner.Run(): %s", err)
	}
	if n != 1 {
		t.Fatalf("Runner.Run(): n = %d", n)
	}
}

func TestRunner_MustRun(t *testing.T) {
	r := New()
	if r.MustRun(NewTaskFromFunc(nil)) == nil {
		t.Fatal("Runner.MustRun(): return nil")
	}
}

func TestRunner_Exit(t *testing.T) {
	r := New()
	r.MustRun(NewTaskFromFunc(nil))
	if r.Exited() {
		t.Fatal("Runner.Exited(): true")
	}
	if err := r.Exit(); err != nil {
		t.Fatalf("Runner.Exit(): [1] %s", err)
	}
	if !r.Exited() {
		t.Fatal("Runner.Exited(): false")
	}

	if err := r.Run(NewTaskFromFunc(nil)); err != ErrExited {
		t.Fatalf("Runner.Run(): [exited] %s", err)
	}
	if err := r.Exit(); err != nil {
		t.Fatalf("Runner.Exit(): [2] %s", err)
	}

	r = New()
	r.MustRun(NewTaskFromFunc(nil, func() error {
		return errors.New("err1")
	}))
	r.MustRun(NewTaskFromFunc(nil, func() error {
		return errors.New("err2")
	}))
	if err := r.Exit(); err == nil {
		t.Fatal("Runner.Exit(): [3] nil error")
	} else {
		s := err.Error()
		if s != "err2; err1" {
			t.Fatalf("Runner.Exit(): [3] %s", s)
		}
	}

	r = New()
	n := 0
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := r.WaitBy(nil); err == nil {
			n = 1
		} else {
			t.Logf("Runner.WaitBy(): %s", err)
		}
	}()
	time.Sleep(time.Millisecond * 20)
	if err := r.Exit(); err != nil {
		t.Fatalf("Runner.Exit(): [4] %s", err)
	}
	wg.Wait()
	if n != 1 {
		t.Fatalf("Runner.WaitBy(): %d", n)
	}
}

func TestRunner_Wait(t *testing.T) {
	r := New()
	r.MustRun(NewTaskFromFunc(nil, func() error {
		return errors.New("test")
	}))
	var n int
	var err error

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		defer wg.Done()
		err = r.Wait()
	}()
	go func() {
		defer wg.Done()
		WaitSystemExit()
		n = 1
	}()

	time.Sleep(time.Millisecond * 20)
	// Close system exit channel.
	systemExitOnce.Do(func() { close(systemExitChan) })
	wg.Wait()

	if n != 1 {
		t.Fatalf("WaitSystemExit(): %d", n)
	}
	if err == nil {
		t.Fatal("Runner.Wait(): nil error")
	}
}
