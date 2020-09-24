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
	"strings"
)

// Errors defines a collection of errors, which is usually used to represent
// errors of the same type within a cycle within an application.
type Errors struct {
	errs []error
}

// Error method is an implementation of the error interface.
func (e *Errors) Error() string {
	switch l := len(e.errs); l {
	case 0:
		return "<empty errors>"
	case 1:
		return e.errs[0].Error()
	default:
		s := make([]string, 0, l)
		for i := 0; i < l; i++ {
			s = append(s, e.errs[i].Error())
		}
		return strings.Join(s, "; ")
	}
}

// Add method adds an error to the current error list.
// If the given error is nil, it is automatically ignored.
func (e *Errors) Add(err error) {
	if err == nil {
		return
	}
	if v, ok := err.(*Errors); ok {
		if len(v.errs) > 0 {
			e.errs = append(e.errs, v.errs...)
		}
	} else {
		e.errs = append(e.errs, err)
	}
}

// First method returns the first error in the current error list
// or nil if the list is empty.
func (e *Errors) First() error {
	if len(e.errs) == 0 {
		return nil
	}
	return e.errs[0]
}

// Last method returns the last error in the current error list
// or nil if the list is empty.
func (e *Errors) Last() error {
	if len(e.errs) == 0 {
		return nil
	}
	return e.errs[len(e.errs)-1]
}

// Len method returns the length of the current error list.
func (e *Errors) Len() int {
	return len(e.errs)
}

// All method returns all errors in the current error list.
func (e *Errors) All() []error {
	return e.errs
}
