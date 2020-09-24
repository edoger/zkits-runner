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
	"testing"
)

func TestIsPanicError(t *testing.T) {
	items := []struct {
		Err  error
		Want bool
	}{
		{nil, false},
		{errors.New("test"), false},
		{new(PanicError), true},
	}

	for i, item := range items {
		if got := IsPanicError(item.Err); got != item.Want {
			t.Fatalf("IsPanicError(): [%d] %v", i, got)
		}
	}
}

func TestSafeCall(t *testing.T) {
	err1 := SafeCall(func() error { return nil })
	if err1 != nil {
		t.Fatalf("SafeCall(): %s", err1)
	}

	err2 := errors.New("err2")
	err3 := SafeCall(func() error { return err2 })
	if err3 != err2 {
		t.Fatalf("SafeCall(): %s", err3)
	}

	err4 := SafeCall(func() error { panic("err4") })
	if err4 == nil {
		t.Fatal("SafeCall(): nil error")
	}
	if !IsPanicError(err4) {
		t.Fatalf("SafeCall(): %s", err4)
	}
}

func TestMustCall(t *testing.T) {
	var v interface{}

	f1 := func() {
		defer func() { v = recover() }()
		MustCall(func() error { return nil })
	}
	f1()
	if v != nil {
		t.Fatalf("MustCall(): %v", v)
	}

	f2 := func() {
		defer func() { v = recover() }()
		MustCall(func() error { return errors.New("test") })
	}
	f2()
	if v == nil {
		t.Fatal("MustCall(): no panic")
	}

	f3 := func() {
		defer func() { v = recover() }()
		MustCall(func() error {
			return SafeCall(func() error {
				panic("test")
			})
		})
	}
	v = nil
	f3()
	if v == nil {
		t.Fatal("MustCall(): no panic")
	}
}

type testFmtStringerForPanicError string

func (s testFmtStringerForPanicError) String() string {
	return string(s)
}

func TestPanicError_Error(t *testing.T) {
	items := []struct {
		Err  *PanicError
		Want string
	}{
		{&PanicError{"test1"}, "test1"},
		{&PanicError{errors.New("test2")}, "test2"},
		{&PanicError{testFmtStringerForPanicError("test3")}, "test3"},
		{&PanicError{4}, "panic: 4"},
	}

	for i, item := range items {
		if s := item.Err.Error(); s != item.Want {
			t.Fatalf("PanicError.Error(): [%d] %s", i, s)
		}
	}
}
