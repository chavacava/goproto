# goproto
Utility to generate funcs prototypes from Go source code.

Useful when distributing Go applications binary.

```
Usage goproto:
  -i string
    	input file (default stdin)
  -include-comp-comment
    	includes a //go:binary-only-package compilation comment (default true)
  -o string
    	output file (default stdout)
  -public-only
    	only generates prototypes of public functions (default true)
```
