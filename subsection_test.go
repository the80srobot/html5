package html5

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/the80srobot/html5/bindings"
	"github.com/the80srobot/html5/safe"
)

func TestSubsection(t *testing.T) {
	for _, tc := range []struct {
		comment string
		input   *SubsectionNode
		opts    *CompileOptions
		values  []bindings.BindArg
		output  string
	}{
		{
			comment: "two rows",
			input: &SubsectionNode{
				Name: "user_comments",
				Prototype: &MultiNode{
					Contents: []Node{
						&TextNode{Value: bindings.Declare("user", safe.Default)},
						&TextNode{Value: safe.Const(" says ")},
						&TextNode{Value: bindings.Declare("comment", safe.Default)},
						&TextNode{Value: safe.Const("\n")},
					},
				},
			},
			opts: &Debug,
			values: []bindings.BindArg{
				{
					Name: "user_comments",
					NestedRows: [][]bindings.BindArg{
						{{Name: "user", Value: safe.Const("John")}, {Name: "comment", Value: safe.Const("Hello!")}},
						{{Name: "user", Value: safe.Const("Jane")}, {Name: "comment", Value: safe.Const("Howdy!")}},
					},
				},
			},
			output: "John says Hello!\nJane says Howdy!\n",
		},
		{
			comment: "nested",
			input: &SubsectionNode{
				Name: "user_comments",
				Prototype: &MultiNode{
					Contents: []Node{
						&TextNode{Value: bindings.Declare("user", safe.Default)},
						&TextNode{Value: safe.Const(" says ")},
						&TextNode{Value: bindings.Declare("comment", safe.Default)},
						&SubsectionNode{
							Name:      "comment_replies",
							Prototype: &TextNode{Value: bindings.Declare("reply", safe.Default)},
						},
						&TextNode{Value: safe.Const("\n")},
					},
				},
			},
			opts: &Debug,
			values: []bindings.BindArg{
				{
					Name: "user_comments",
					NestedRows: [][]bindings.BindArg{
						{
							{Name: "user", Value: safe.Const("John")},
							{Name: "comment", Value: safe.Const("Hello!")},
							{
								Name: "comment_replies",
								NestedRows: [][]bindings.BindArg{
									{{Name: "reply", Value: safe.Const("Love!")}},
									{{Name: "reply", Value: safe.Const("Good to see you!")}},
								},
							},
						},
						{{Name: "user", Value: safe.Const("Jane")}, {Name: "comment", Value: safe.Const("Howdy!")}},
					},
				},
			},
			output: "John says Hello!Love!Good to see you!\nJane says Howdy!\n",
		},
	} {
		t.Run(tc.comment, func(t *testing.T) {
			if diff := cmp.Diff(tc.output, mustGenerateHTML(t, tc.input, tc.opts, tc.values)); diff != "" {
				t.Errorf("GenerateHTML(%v, %v, %v)\n => (-)wanted vs (+)got:\n%s", tc.input, tc.opts, tc.values, diff)
			}
		})
	}
}
