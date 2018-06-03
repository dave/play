package builderjs

import (
	"context"
	"go/token"
	"go/types"
	"sort"

	"go/ast"
	"go/parser"

	"fmt"

	"strings"

	"bytes"
	"crypto/sha1"

	"github.com/dave/services/includer"
	"github.com/gopherjs/gopherjs/compiler"
	"golang.org/x/tools/go/gcexportdata"
)

func BuildPackage(path string, source map[string]map[string]string, tags []string, deps []*compiler.Archive, minify bool, archives map[string]*compiler.Archive, packages map[string]*types.Package) (*compiler.Archive, error) {

	for _, a := range deps {
		if archives[a.ImportPath] == nil {
			archives[a.ImportPath] = a
		}
		if packages[a.ImportPath] == nil {
			p, err := gcexportdata.Read(bytes.NewReader(a.ExportData), token.NewFileSet(), packages, a.ImportPath)
			if err != nil {
				return nil, err
			}
			packages[a.ImportPath] = p
		}
	}

	fset := token.NewFileSet()

	var importContext *compiler.ImportContext
	importContext = &compiler.ImportContext{
		Packages: packages,
		Import: func(imp string) (*compiler.Archive, error) {
			a, ok := archives[imp]
			if ok {
				return a, nil
			}
			sourceFiles, ok := source[imp]
			if ok {
				// We have the source for this dep
				archive, err := compileFiles(fset, imp, tags, sourceFiles, importContext, minify)
				if err != nil {
					return nil, err
				}
				return archive, nil
			}
			return nil, fmt.Errorf("%s not found", imp)
		},
	}

	archive, err := importContext.Import(path)
	if err != nil {
		return nil, err
	}

	return archive, nil
}

func compileFiles(fset *token.FileSet, path string, tags []string, sourceFiles map[string]string, importContext *compiler.ImportContext, minify bool) (*compiler.Archive, error) {
	var files []*ast.File
	inc := includer.New(sourceFiles, tags)
	for name, contents := range sourceFiles {
		include, err := inc.Include(name)
		if err != nil {
			return nil, err
		}
		if !include {
			continue
		}
		f, err := parser.ParseFile(fset, name, contents, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no buildable Go source files in %s", path)
	}

	// TODO: Remove this when https://github.com/gopherjs/gopherjs/pull/742 is merged
	// Files must be in the same order to get reproducible JS
	sort.Slice(files, func(i, j int) bool {
		return fset.File(files[i].Pos()).Name() > fset.File(files[j].Pos()).Name()
	})

	archive, err := compiler.Compile(path, files, fset, importContext, minify)
	if err != nil {
		return nil, err
	}

	for name, contents := range sourceFiles {
		if !strings.HasSuffix(name, ".inc.js") {
			continue
		}
		archive.IncJSCode = append(archive.IncJSCode, []byte("\t(function() {\n")...)
		archive.IncJSCode = append(archive.IncJSCode, []byte(contents)...)
		archive.IncJSCode = append(archive.IncJSCode, []byte("\n\t}).call($global);\n")...)
	}

	return archive, nil
}

func GetPackageCode(ctx context.Context, archive *compiler.Archive, minify, initializer bool) (contents []byte, hash []byte, err error) {
	dceSelection := make(map[*compiler.Decl]struct{})
	for _, d := range archive.Declarations {
		dceSelection[d] = struct{}{}
	}
	buf := new(bytes.Buffer)

	if initializer {
		var s string
		if minify {
			s = `$load["%s"]=function(){`
		} else {
			s = `$load["%s"] = function () {` + "\n"
		}
		if _, err := fmt.Fprintf(buf, s, archive.ImportPath); err != nil {
			return nil, nil, err
		}
	}

	if err := compiler.WritePkgCode(archive, dceSelection, minify, &compiler.SourceMapFilter{Writer: buf}); err != nil {
		return nil, nil, err
	}

	if minify {
		// compiler.WritePkgCode always finishes with a "\n". In minified mode we should remove this.
		buf.Truncate(buf.Len() - 1)
	}

	if initializer {
		/*
			var s string
			if minify {
				s = "};$done();"
			} else {
				s = "};\n$done();"
			}
			if _, err := fmt.Fprint(buf, s); err != nil {
				return nil, nil, err
			}
		*/
		if _, err := fmt.Fprint(buf, "};"); err != nil {
			return nil, nil, err
		}
	}

	sha := sha1.New()
	if _, err := sha.Write(buf.Bytes()); err != nil {
		return nil, nil, err
	}
	return buf.Bytes(), sha.Sum(nil), nil
}
