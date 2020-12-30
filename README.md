# Go package for generating HTML5

IMPORTANT: THIS IS STILL EXPERIMENTAL AND SUBJECT TO CHANGE

This package provides a template-free, declarative mechanism for generating
HTML5. It has some advantages over `html/template`:

* It is about 5-10 times faster
* With careful use, allows streaming HTML without holding the entire page in
  memory
* Can generate tidy output and minified output directly, without needing
  additional passes
* Makes it easy to build reusable, composable pieces
* HTML is generated from pure go code, which can sometimes improve readability
  and locality, especially in simple applications

## Safety

Many HTML libraries divide the components of an HTML page into two categories:

1. Trusted, static elements (the template source itself)
2. Untrusted, dynamic elements (string values interpolated at runtime)

The above threat model is too coarse, resulting in low runtime performance and
unnecessary implementation complexity.

Instead, we introduce a third category:

3. Dynamic elements obtained from a trusted source or method of construction

Examples of this third category include:

* Build-time constants
* Strings generated from numbers, safe data structures, etc
* Data loaded from the server's disk (assuming sane security properties)

This package allows the programmer to insert trusted strings into various
contexts, but requires explicit trust annotations. Fully trusted strings are
allowed in any context, while untrusted strings can sometimes be escaped to be
made fit for use in some contexts. Some contexts always require fully trusted
strings.

## Performance

The runtime performance of this package compares favorably to `html/template`,
being between 5-10 times faster in microbenchmarks.

```
html5 % go test -bench=. ./...
goos: darwin
goarch: amd64
pkg: github.com/the80srobot/html5
BenchmarkSmallTemplate-16    	  379461	      2996 ns/op
BenchmarkSmallPage-16        	 3247833	       391 ns/op
PASS
ok  	github.com/the80srobot/html5	2.889s
PASS
ok  	github.com/the80srobot/html5/html	0.070s
```