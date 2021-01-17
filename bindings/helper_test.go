package bindings

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/the80srobot/html5/safe"
)

func logDebug(t testing.TB, dd DebugDumper, msg string) {
	t.Helper()

	var sb strings.Builder
	dd.DebugDump(&sb, 0)
	t.Logf("%s: %s", msg, sb.String())
}

// Tests that Bind can construct an entire structure of nested Maps correctly
// and implicitly, given only the names.
func TestBindNames(t *testing.T) {
	args := BindArgs{
		{Name: "title", Value: safe.Const("Articles")},
		{
			Name: "articles",
			NestedRows: [][]BindArg{
				{
					{Name: "title", Value: safe.Const("Hello, World!")},
					{Name: "author", Value: safe.Const("Adam")},
					{
						Name: "tags",
						NestedRows: [][]BindArg{
							{{Name: "tag", Value: safe.Const("diary")}},
							{{Name: "tag", Value: safe.Const("blog")}},
						},
					},
				},
				{
					{Name: "title", Value: safe.Const("Recipe for Soup")},
					{Name: "author", Value: safe.Const("Link's Grandma")},
					{
						Name: "tags",
						NestedRows: [][]BindArg{
							{{Name: "tag", Value: safe.Const("cooking")}},
							{{Name: "tag", Value: safe.Const("recipes")}},
						},
					},
				},
			},
		},
	}

	logDebug(t, args, "Bind arguments")
	// Create a ValueMap using the Bind helper.
	var m Map
	got := m.MustBind()
	if err := Bind(got, args...); err != nil {
		t.Fatal(err)
	}

	logDebug(t, got, "Map from Bind")

	// Create a ValueMap manually. The two should match.
	m = Map{}
	title := m.Declare("title", safe.Default)
	articles := m.Nest("articles")
	articleTitle := articles.Declare("title", safe.Default)
	author := articles.Declare("author", safe.Default)
	tags := articles.Nest("tags")
	tag := tags.Declare("tag", safe.Default)

	a1 := articles.MustBind(
		articleTitle.BindConst("Hello, World!"),
		author.BindConst("Adam"),
		tags.BindSeries(
			tags.MustBind(tag.BindConst("diary")),
			tags.MustBind(tag.BindConst("blog"))))

	a2 := articles.MustBind(
		articleTitle.BindConst("Recipe for Soup"),
		author.BindConst("Link's Grandma"),
		tags.BindSeries(
			tags.MustBind(tag.BindConst("cooking")),
			tags.MustBind(tag.BindConst("recipes"))))

	want := m.MustBind(title.BindConst("Articles"), articles.BindSeries(a1, a2))

	opt := cmp.Transformer("ValueMap", func(vm *ValueMap) string {
		return DebugString(vm)
	})

	logDebug(t, want, "Map populated manually")

	if diff := cmp.Diff(got, want, opt); diff != "" {
		t.Errorf("Bind using the logged args => (-)wanted vs (+)got:\n%s", diff)
	}
}
