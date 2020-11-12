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
	"strings"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	if r := New(); r == nil {
		t.Fatal("New(): nil")
	}
}

func TestRunner_Run(t *testing.T) {
	do := func(fs ...func(Runner)) {
		for _, f := range fs {
			f(New())
		}
	}

	do(func(r Runner) {
		if got := r.Run(NewTaskFromFunc(nil)); got != nil {
			t.Fatalf("Runner.Run(): %s", got)
		}
	}, func(r Runner) {
		want := errors.New("test error")
		if got := r.Run(NewTaskFromFunc(func() error { return want })); got != want {
			t.Fatalf("Runner.Run(): %s", got)
		}
	}, func(r Runner) {
		if got := r.Run(NewTaskFromFunc(func() error { panic("test") })); !IsPanicError(got) {
			t.Fatalf("Runner.Run(): %s", got)
		}
	}, func(r Runner) {
		if err := r.Exit(); err != nil {
			t.Fatal(err)
		}

		if got := r.Run(NewTaskFromFunc(nil)); got != ErrExited {
			t.Fatalf("Runner.Run(): %s", got)
		}
	})
}

func TestRunner_MustRun(t *testing.T) {
	do := func(fs ...func(Runner)) {
		for _, f := range fs {
			f(New())
		}
	}

	do(func(r Runner) {
		r.MustRun(NewTaskFromFunc(nil))
	}, func(r Runner) {
		defer func() {
			if recover() == nil {
				t.Fatal("Runner.MustRun(): no panic")
			}
		}()

		r.MustRun(NewTaskFromFunc(func() error {
			return errors.New("test")
		}))
	}, func(r Runner) {
		if err := r.Exit(); err != nil {
			t.Fatal(err)
		}

		defer func() {
			if recover() == nil {
				t.Fatal("Runner.MustRun(): no panic")
			}
		}()

		r.MustRun(NewTaskFromFunc(nil))
	})
}

func TestRunner_Exit(t *testing.T) {
	do := func(fs ...func(Runner)) {
		for _, f := range fs {
			f(New())
		}
	}

	do(func(r Runner) {
		var ss []string
		r.MustRun(NewTaskFromFunc(func() error {
			ss = append(ss, "A1")
			return nil
		}, func() error {
			ss = append(ss, "B1")
			return nil
		}))
		r.MustRun(NewTaskFromFunc(func() error {
			ss = append(ss, "A2")
			return nil
		}, func() error {
			ss = append(ss, "B2")
			return nil
		}))

		if got := r.Exited(); got {
			t.Fatal("Runner.Exited(): true")
		}
		if err := r.Exit(); err != nil {
			t.Fatalf("Runner.Exit(): %s", err)
		}
		if got := r.Exited(); !got {
			t.Fatal("Runner.Exited(): false")
		}

		if got := strings.Join(ss, "-"); got != "A1-A2-B2-B1" {
			t.Fatalf("Runner.Exit(): %s", got)
		}
	}, func(r Runner) {
		r.MustRun(NewTaskFromFunc(nil, func() error {
			return errors.New("err1")
		}))
		r.MustRun(NewTaskFromFunc(nil, func() error {
			return errors.New("err2")
		}))

		if err := r.Exit(); err == nil {
			t.Fatal("Runner.Exit(): no error")
		} else {
			if got := err.Error(); got != "err2; err1" {
				t.Fatalf("Runner.Exit(): %s", got)
			}
		}
	})
}

func TestRunner_Wait(t *testing.T) {
	do := func(fs ...func(Runner)) {
		for _, f := range fs {
			f(New())
		}
	}

	do(func(r Runner) {
		r.MustRun(NewTaskFromFunc(nil))

		var m, n int

		ch := make(chan struct{})
		wg := new(sync.WaitGroup)
		wg.Add(2)
		go func() {
			defer wg.Done()
			if err := r.Wait(); err != nil {
				t.Fatalf("Runner.Wait(): %s", err)
			} else {
				m++
			}
		}()
		go func() {
			defer wg.Done()
			if err := r.WaitBy(ch); err != nil {
				t.Fatalf("Runner.WaitBy(): %s", err)
			} else {
				n++
			}
		}()

		close(ch)
		wg.Wait()

		if m != 1 || n != 1 {
			t.Fatalf("Runner: %d %d", m, n)
		}
	})
}
