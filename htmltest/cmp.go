package htmltest

import (
	"github.com/google/go-cmp/cmp"
	"github.com/the80srobot/html5/html"
)

var CmpOpts = []cmp.Option{
	cmp.Comparer(func(x, y html.SafeString) bool { return x == y }),
}
