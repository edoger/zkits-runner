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

func TestErrors(t *testing.T) {
	errs := new(Errors)

	if s := errs.Error(); s != "<empty errors>" {
		t.Fatalf("Errors.Error(): %s", s)
	}
	if err := errs.First(); err != nil {
		t.Fatalf("Errors.First(): %s", err)
	}
	if err := errs.Last(); err != nil {
		t.Fatalf("Errors.Last(): %s", err)
	}
	if n := errs.Len(); n != 0 {
		t.Fatalf("Errors.Last(): %d", n)
	}
	if v := errs.All(); len(v) != 0 {
		t.Fatalf("Errors.All(): %v", v)
	}

	errs.Add(nil)
	if s := errs.Error(); s != "<empty errors>" {
		t.Fatalf("Errors.Error(): %s", s)
	}

	errs.Add(errors.New("test1"))
	if s := errs.Error(); s != "test1" {
		t.Fatalf("Errors.Error(): %s", s)
	}

	errs.Add(errors.New("test2"))
	if s := errs.Error(); s != "test1; test2" {
		t.Fatalf("Errors.Error(): %s", s)
	}

	errs2 := new(Errors)
	errs2.Add(errors.New("test3"))
	errs.Add(errs2)
	if s := errs.Error(); s != "test1; test2; test3" {
		t.Fatalf("Errors.Error(): %s", s)
	}

	if err := errs.First(); err == nil {
		t.Fatal("Errors.First(): nil error")
	}
	if err := errs.Last(); err == nil {
		t.Fatal("Errors.Last(): nil error")
	}
	if n := errs.Len(); n != 3 {
		t.Fatalf("Errors.Last(): %d", n)
	}
	if v := errs.All(); len(v) != 3 {
		t.Fatalf("Errors.All(): %v", v)
	}
}
