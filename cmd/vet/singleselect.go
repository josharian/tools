// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
This file contains code to check for select statements that contain only a
single case, but no default: case. Such a statement is useless, because the
select can be removed, and confusing because it looks like a non-blocking
channel operation. In many cases, authors may be expecting non-blocking
semantics, but forget the empty default, causing hard-to-debug errors or
deadlocks.

For example:

	select {
	case ch <- true:
		// code
	}

Should be rewritten as:

	ch <- true
	// code

if it is meant to be blocking, and as

	select {
	case ch <- true:
		// code
	default:
	}

if it is meant to be non-blocking.
*/

package main

import "go/ast"

func init() {
	register("singleselect",
		"check for blocking single-case select statements",
		checkSelectCases,
		selectStmt)
}

// checkSelectCases enumerates the cases of the provided select statement,
// checking if there is only one, and if so, whether there is a default: case
// statement to enable non-blocking behavior.
func checkSelectCases(f *File, node ast.Node) {
	n := node.(*ast.SelectStmt)

	defaults := false
	cases := 0
	for _, condition := range n.Body.List {
		c := condition.(*ast.CommClause)

		switch c.Comm.(type) {
		case *ast.SendStmt:
			// chan<-
		case *ast.ExprStmt:
			// <-chan
		case nil:
			// select has default case
			defaults = true
		}
		cases++
	}

	if cases != 1 {
		// select is not useless
		// is either blocking select {}
		// or multi-case
		return
	}

	if defaults {
		f.Bad(n.Pos(), "select is useless with only default case")
	} else {
		f.Bad(n.Pos(), "select has only single case and no default")
	}
}
