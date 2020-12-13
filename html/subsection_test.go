package html

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSubsection(t *testing.T) {
	for _, tc := range []struct {
		comment string
		input   *SubsectionNode
		opts    *CompileOptions
		values  []ValueArg
		depth   int
		output  string
	}{
		{
			comment: "two rows",
			input: &SubsectionNode{
				Name: "user_comments",
				Prototype: &MultiNode{
					Contents: []Node{
						&TextNode{Value: Binding("user")},
						&TextNode{Value: FullyTrustedString(" says ")},
						&TextNode{Value: Binding("comment")},
						&TextNode{Value: FullyTrustedString("\n")},
					},
				},
			},
			opts: &Debug,
			values: []ValueArg{
				{
					Name: "user_comments",
					Subsections: [][]ValueArg{
						{{Name: "user", SafeString: FullyTrustedString("John")}, {Name: "comment", SafeString: FullyTrustedString("Hello!")}},
						{{Name: "user", SafeString: FullyTrustedString("Jane")}, {Name: "comment", SafeString: FullyTrustedString("Howdy!")}},
					},
				},
			},
			depth:  0,
			output: "John says Hello!\nJane says Howdy!\n",
		},
		{
			comment: "nested",
			input: &SubsectionNode{
				Name: "user_comments",
				Prototype: &MultiNode{
					Contents: []Node{
						&TextNode{Value: Binding("user")},
						&TextNode{Value: FullyTrustedString(" says ")},
						&TextNode{Value: Binding("comment")},
						&SubsectionNode{
							Name:      "comment_replies",
							Prototype: &TextNode{Value: Binding("reply")},
						},
						&TextNode{Value: FullyTrustedString("\n")},
					},
				},
			},
			opts: &Debug,
			values: []ValueArg{
				{
					Name: "user_comments",
					Subsections: [][]ValueArg{
						{
							{Name: "user", SafeString: FullyTrustedString("John")},
							{Name: "comment", SafeString: FullyTrustedString("Hello!")},
							{
								Name: "comment_replies",
								Subsections: [][]ValueArg{
									{{Name: "reply", SafeString: FullyTrustedString("Love!")}},
									{{Name: "reply", SafeString: FullyTrustedString("Good to see you!")}},
								},
							},
						},
						{{Name: "user", SafeString: FullyTrustedString("Jane")}, {Name: "comment", SafeString: FullyTrustedString("Howdy!")}},
					},
				},
			},
			depth:  0,
			output: "John says Hello!Love!Good to see you!\nJane says Howdy!\n",
		},
	} {
		t.Run(tc.comment, func(t *testing.T) {
			if diff := cmp.Diff(tc.output, mustGenerateHTML(t, tc.input, tc.depth, tc.opts, tc.values)); diff != "" {
				t.Errorf("GenerateHTML(%v, %v, %v, %v)\n => (-)wanted vs (+)got:\n%s", tc.input, tc.depth, tc.opts, tc.values, diff)
			}
		})
	}
}
