// Copyright 2018 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package pretty_test

import (
	"fmt"
	"strings"

	"github.com/cockroachdb/cockroach/pkg/util/pretty"
)

// ExampleTree demonstrates the Tree example from the paper.
func Example_tree() {
	type Tree struct {
		s string
		n []Tree
	}
	tree := Tree{
		"aaa",
		[]Tree{
			{
				"bbbbb",
				[]Tree{
					{s: "ccc"},
					{s: "dd"},
				},
			},
			{s: "eee"},
			{
				"ffff",
				[]Tree{
					{s: "gg"},
					{s: "hhh"},
					{s: "ii"},
				},
			},
		},
	}
	var (
		showTree    func(Tree) pretty.Doc
		showTrees   func([]Tree) pretty.Doc
		showBracket func([]Tree) pretty.Doc
	)
	showTrees = func(ts []Tree) pretty.Doc {
		if len(ts) == 1 {
			return showTree(ts[0])
		}
		return pretty.Fold(pretty.Concat,
			showTree(ts[0]),
			pretty.Text(","),
			pretty.Line,
			showTrees(ts[1:]),
		)
	}
	showBracket = func(ts []Tree) pretty.Doc {
		if len(ts) == 0 {
			return pretty.Nil
		}
		return pretty.Fold(pretty.Concat,
			pretty.Text("["),
			pretty.Nest(1, " ", showTrees(ts)),
			pretty.Text("]"),
		)
	}
	showTree = func(t Tree) pretty.Doc {
		return pretty.Group(pretty.Concat(
			pretty.Text(t.s),
			pretty.Nest(len(t.s), strings.Repeat(" ", len(t.s)), showBracket(t.n)),
		))
	}
	for _, n := range []int{1, 30, 80} {
		p := pretty.Pretty(showTree(tree), n)
		fmt.Printf("%d:\n%s\n\n", n, p)
	}
	// Output:
	// 1:
	// aaa[bbbbb[ccc,
	//           dd],
	//     eee,
	//     ffff[gg,
	//          hhh,
	//          ii]]
	//
	// 30:
	// aaa[bbbbb[ccc, dd],
	//     eee,
	//     ffff[gg, hhh, ii]]
	//
	// 80:
	// aaa[bbbbb[ccc, dd], eee, ffff[gg, hhh, ii]]
}
