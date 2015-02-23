// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains tests for the singleselect checker.

package testdata

func SingleSelectTests() {
	select {
	case <-ch1:
	default:
	}
	select {
	case ch2 <- true:
	default:
	}
	select { // ERROR "select has only single case and no default"
	case <-ch1:
	}
	select { // ERROR "select has only single case and no default"
	case ch2 <- true:
	}
	select {}
}
