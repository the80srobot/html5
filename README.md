# Go package for generating HTML5

This package provides a template-free, declarative mechanism for generating
HTML5. It has some advantages over `template/html`:

* With careful use, allows streaming HTML without holding the entire page in
  memory
* Automatically indents and justifies the output
* Makes it easy to build reusable, composable pieces
* HTML is generated from pure go code, which can sometimes improve readability
  and locality, especially in simple applications

  See the godoc for more information.