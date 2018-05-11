package execute

import (
	"bufio"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
)

const (
	requiresNonil = iota
	requiresExp   = iota
)

func Execute() {
	srcFile := flag.String("i", "", "input file (defaults to stdin)")
	outFile := flag.String("o", "", "output file (defaults to stdout)")
	onlyPublic := flag.Bool("public-only", true, "only generates prototypes of public functions")
	binaryHeader := flag.Bool("include-comp-comment", true, "includes a //go:binary-only-package compilation comment")
	flag.Parse()

	src, err := ioutil.ReadFile(*srcFile)
	if err != nil {
		log.Fatalf("could not open input file: %v", err)
	}

	var w io.Writer
	if *outFile != "" {
		outFile, err := os.Create(*outFile)
		if err != nil {
			log.Fatal(err)
		}
		defer outFile.Close()
		w = bufio.NewWriter(outFile)
	} else {
		w = os.Stdout
	}

	if *binaryHeader {
		fmt.Fprintf(w, "//go:binary-only-package\n\n")
	}

	err = analyzeCode(src, *srcFile, *onlyPublic, w)

	if err != nil {
		log.Fatalf("could not analyze source code: %v", err)
	}
}

func analyzeCode(src []byte, fileName string, onlyPublic bool, w io.Writer) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, fileName, src, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("could not parse input code: %v", err)
	}

	if err := writePackageDeclaration(w, file); err != nil {
		return fmt.Errorf("could not generate package declaration: %v", err)
	}

	ast.Inspect(file, func(x ast.Node) bool {
		switch n := x.(type) {
		case *ast.FuncDecl:
			if !n.Name.IsExported() && onlyPublic { // skip private functions
				return true
			}

			writeFuncPrototype(w, n, fset)
		default:
			return true
		}

		return false
	})

	return nil
}

func writePackageDeclaration(w io.Writer, file *ast.File) error {
	fmt.Fprintf(w, "/*\n%s*/\n", file.Doc.Text())
	_, err := fmt.Fprintf(w, "package %s\n\n", file.Name.String())

	return err

}

func writeFuncPrototype(w io.Writer, fd *ast.FuncDecl, fset *token.FileSet) error {
	if fd == nil {
		log.Fatal("can not extract prototype of a nil function")
	}

	fd.Body = nil // function body amputation before printing

	err := printer.Fprint(w, fset, fd)
	if err != nil {
		log.Fatalf("error while printing prototype of function %s", fd.Name.String())
	}
	fmt.Fprintf(w, "\n\n")

	return nil
}
