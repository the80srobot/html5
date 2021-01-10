package bindings

import (
	"testing"

	"github.com/the80srobot/html5/safe"
)

type article struct {
	title, author, text safe.String
	comments            []comment
	tags                []safe.String
}

type comment struct {
	author, text safe.String
}

var articleData = []article{
	{
		title:  safe.Const("The Pros and Cons of Brushing Your Teeth with Battery Acid"),
		author: safe.Const("Alice"),
		text:   safe.Const("Pros: it's metal, Cons: none I can think of"),
		comments: []comment{
			{
				author: safe.Const("Bob"),
				text:   safe.Const("I've been brushing my teeth with battery acid for years now, and I'm starting to suspect it's indeed very metal."),
			},
			{
				author: safe.Const("Your Dentist"),
				text:   safe.Const("I endorse this message"),
			},
		},
		tags: []safe.String{safe.Const("dentistry"), safe.Const("metal"), safe.Const("deep thoughts")},
	},
	{
		title:  safe.Const("Pictures of My Cat"),
		author: safe.Const("Bob"),
		text:   safe.Const("So fluffy"),
	},
}

func TestNestedMapsRoundTrip(t *testing.T) {
	// Declare the binding structure. Top level has a stream of articles, and
	// each article has a stream of comments and a stream of tags.
	var bindings Map

	pageTitle := bindings.Declare("page_title", safe.HTMLSafe)

	articleBindings := bindings.Nest("articles")
	title := articleBindings.Declare("title", safe.HTMLSafe)
	author := articleBindings.Declare("author", safe.HTMLSafe)
	text := articleBindings.Declare("text", safe.HTMLSafe)

	commentBindings := articleBindings.Nest("comments")
	commentAuthor := commentBindings.Declare("author", safe.TextSafe)
	commentText := commentBindings.Declare("text", safe.TextSafe)

	tagBindings := articleBindings.Nest("tags")
	tag := tagBindings.Declare("tag", safe.TextSafe)

	t.Logf("Article, comment and tag bindings: %v", &bindings)

	// Instantiate a value map and populate it with articles and comments under
	// the articles.
	values, err := bindings.Bind(pageTitle.BindConst("Articles"))
	if err != nil {
		t.Fatalf("couldn't create the root value map: %v", err)
	}

	var articles ValueSeries
	for _, a := range articleData {
		article, err := articleBindings.Bind(title.Bind(a.title), author.Bind(a.author), text.Bind(a.text))
		if err != nil {
			t.Fatalf("couldn't create an article value map: %v", err)
		}

		var comments ValueSeries
		for _, c := range a.comments {
			comment, err := commentBindings.Bind(commentAuthor.Bind(c.author), commentText.Bind(c.text))
			if err != nil {
				t.Fatalf("couldn't create a comment value map: %v", err)
			}
			comments = append(comments, comment)
		}

		var tags ValueSeries
		for _, tg := range a.tags {
			tag, err := tagBindings.Bind(tag.Bind(tg))
			if err != nil {
				t.Fatalf("couldn't create a tag value map: %v", err)
			}
			tags = append(tags, tag)
		}

		article.Set(commentBindings.BindSeries(comments...))
		article.Set(tagBindings.BindSeries(tags...))
		articles = append(articles, article)
	}

	values.Set(articleBindings.BindSeries(articles...))
	t.Logf("Value dump: %v", values)

	// Check that the data we get back by iterating articles and their comments
	// looks like what we put in above.

	if s := values.GetString(pageTitle); s != "Articles" {
		t.Errorf("GetString(%v) => %q, wanted 'Articles'", title, s)
	}

	iter := values.GetStream(articleBindings).Stream()
	i := 0
	for a := iter(); a != nil; a = iter() {
		t.Logf("article %d: %v", i, a)
		if s := a.GetString(title); s != articleData[i].title.String() {
			t.Errorf("(article) GetString(%v) => %q, wanted %q", title, s, articleData[i].title)
		}
		if s := a.GetString(author); s != articleData[i].author.String() {
			t.Errorf("(article) GetString(%v) => %q, wanted %q", author, s, articleData[i].author)
		}

		iter2 := a.GetStream(commentBindings).Stream()
		j := 0
		for c := iter2(); c != nil; c = iter2() {
			t.Logf("comment %d: %v", j, c)

			if len(articleData[i].comments) <= j {
				t.Fatalf("article %d contains %d comments, but iterator is at %d", i, len(articleData[i].comments), j)
			}

			if s := c.GetString(commentAuthor); s != articleData[i].comments[j].author.String() {
				t.Errorf("(comment) GetString(%v) => %q, wanted %q", commentAuthor, s, articleData[i].comments[j].author)
			}
			j++
		}

		i++
	}
}

func TestNestedMapWrongBindingGet(t *testing.T) {
	var page Map
	pageTitle := page.Declare("title", safe.HTMLSafe)

	articleBindings := page.Nest("articles")
	articleTitle := articleBindings.Declare("title", safe.HTMLSafe)

	pageValues := page.MustBind(pageTitle.BindConst("Welcome to Zombo.com"))

	articleValues := articleBindings.MustBind(articleTitle.BindConst("Anything is possible at Zombo.com"))
	pageValues.Set(articleBindings.BindSeries(articleValues))

	// Trying to use the page-level binding on the article values should fail,
	// and vice-versa. Because this is considered a programmer error, GetString
	// should panic.
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected GetString(%v) to panic", pageTitle)
			r = nil
		}
	}()
	articleValues.GetString(pageTitle)
}

func TestNestedMapWrongBindingSet(t *testing.T) {
	var page Map
	page.Declare("title", safe.HTMLSafe)

	articleBindings := page.Nest("articles")
	articleTitle := articleBindings.Declare("title", safe.HTMLSafe)
	pageValues := page.MustBind()

	// Trying to use the article-level var to set a value on the page-level
	// value map should panic. (Considered a programmer error.)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected Set(%v) to panic", articleTitle)
			r = nil
		}
	}()
	pageValues.Set(articleTitle.BindConst("Welcome to Zombo.com"))
}

func TestVarAttach(t *testing.T) {
	v := Declare("foo")
	var m Map

	v = m.Attach(v, safe.Default)

	// The variable should now be valid for the map. Test that out by using it
	// to bind a value.
	vm, err := m.Bind(v.BindConst("bar"))
	if err != nil {
		t.Fatal(err)
	}

	s := vm.GetString(v)
	if s != "bar" {
		t.Errorf("ValueMap.GetString(%v) => %q, wanted %q", v, s, "bar")
	}
}

func TestTrustClimbing(t *testing.T) {
	var m Map
	v := m.Declare("comment_text", safe.TextSafe)
	if _, err := v.tryBind(safe.UntrustedString("Hello!")); err == nil {
		t.Error("Var.Set() should fail to set an untrusted string on TextSafe Var")
	}

	// Declaring the var again with the default trust shouldn't change anything
	// - text should still be the right level.
	v = m.Declare("comment_text", safe.Default)
	if _, err := v.tryBind(safe.Bless(safe.TextSafe, "Hello!")); err != nil {
		t.Errorf("Var.Set() of a TextSafe string: %v", err)
	}

	// This should climb all the way up to fully trusted, as the only way to
	// reconcile the two trust levels.
	v = m.Declare("comment_text", safe.URLSafe)
	if _, err := v.tryBind(safe.Bless(safe.TextSafe, "Hello!")); err == nil {
		t.Error("Var.Set() should have refused to set TextSafe after trust climbing to fully trusted")
	}

	// Fully trusted strings should still be accepted.
	if _, err := v.tryBind(safe.Bless(safe.FullyTrusted, "Hello!")); err != nil {
		t.Errorf("Var.Set() of a FullyTrusted string: %v", err)
	}
}

// func TestMapString(t *testing.T) {
// 	var m Map
// 	t1 := m.DeclareString("foo", Untrusted)
// 	m.DeclareString("bar", Untrusted)

// 	if t3 := m.DeclareString("foo", AttributeSafe); t3 != t1 {
// 		t.Errorf("Different tag when declaring the same value a second time (%d vs %d)", t3, t1)
// 	}

// 	if s := m.String(t1); s != s1 {
// 		t.Errorf("Lookup after declare yielded the wrong string %q (wanted %q)", s.Name, s1.Name)
// 	}
// }

// func BenchmarkMap100LookupBaseline(b *testing.B) {
// 	m := map[string]Var{}
// 	for i := 0; i < 100; i++ {
// 		k := fmt.Sprintf("key_%d", i)
// 		m[k] = Var{Name: k}
// 	}

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		v := m["key_50"]
// 		// This check is mainly here to make sure the lookup isn't optimized
// 		// away.
// 		if v.Name == "" {
// 			b.Fatal("bad lookup")
// 		}
// 	}
// }

// func BenchmarkMap100Lookup(b *testing.B) {
// 	m := NewMap()
// 	tags := []Tag{}
// 	for i := 0; i < 100; i++ {
// 		k := fmt.Sprintf("key_%d", i)
// 		tags = append(tags, m.DeclareString(Var{Name: k}))
// 	}

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		v := m.String(tags[50])
// 		// This check is here because we also have it in the baseline.
// 		if v.Name == "" {
// 			b.Fatal("bad lookup")
// 		}
// 	}
// }
