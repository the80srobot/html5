package bindings

import (
	"fmt"

	"github.com/the80srobot/html5/safe"
)

func Example() {
	var m Map

	userName := m.Declare("author_name", safe.TextSafe)
	articleURL := m.Declare("article_url", safe.URLSafe)
	articleHTML := m.Declare("article_html", safe.HTMLSafe)

	comments := m.Nest("comments")
	commentHTML := comments.Declare("comment_html", safe.HTMLSafe)
	commentAuthor := comments.Declare("author", safe.TextSafe)

	values := m.MustBind(
		userName.BindConst("adam"),
		articleURL.BindConst("https://something.com"),
		articleHTML.BindConst("<p>...</p>"),
		comments.BindSeries(
			comments.MustBind(
				commentHTML.BindConst("<p>Hello!</p>"),
				commentAuthor.BindConst("Peter")),
			comments.MustBind(
				commentHTML.BindConst("<p>Hello!</p>"),
				commentAuthor.BindConst("Paul"))))

	commentStream := values.GetStream(comments)
	firstCommentValue := commentStream.Stream()()
	fmt.Printf("First comment's author: %s", firstCommentValue.GetString(commentAuthor))

	// Output: First comment's author: Peter
}
