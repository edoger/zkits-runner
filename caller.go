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
	"fmt"
)

// PanicError defines the panic error captured by recover.
// We do not recommend using this error type in application business logic.
// The purpose of designing this error type is to ensure that the SafeCall
// function can report panic.
type PanicError struct {
	v interface{}
}

// Error method is an implementation of the error interface.
func (e *PanicError) Error() string {
	switch o := e.v.(type) {
	case string:
		return o
	case error:
		return o.Error()
	case fmt.Stringer:
		return o.String()
	}
	return fmt.Sprintf("panic: %v", e.v)
}

// IsPanicError determines whether the given error is a PanicError.
// If the given error is nil, it always returns false.
func IsPanicError(err error) (ok bool) {
	if err != nil {
		_, ok = err.(*PanicError)
	}
	return
}

// SafeCall executes the given function immediately, if a panic occurs,
// it returns a PanicError, otherwise it returns the error returned
// by the given function.
func SafeCall(f func() error) (err error) {
	defer func() {
		if v := recover(); v != nil {
			err = &PanicError{v: v}
		}
	}()
	err = f()
	return
}

// MustCall executes the given function immediately, and panic immediately
// if the given function returns a non-nil error.
func MustCall(f func() error) {
	if err := f(); err != nil {
		// If this error is already a PanicError, we don't need to
		// wrap the underlying panic again.
		if e, ok := err.(*PanicError); ok {
			panic(e.v)
		}
		panic(err)
	}
}
