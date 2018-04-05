package stores

import (
	"go/parser"
	"go/token"

	"strconv"

	"sort"

	"strings"

	"github.com/dave/flux"
	"github.com/dave/play/actions"
)

func NewScannerStore(app *App) *ScannerStore {
	s := &ScannerStore{
		app:     app,
		imports: map[string]map[string][]string{},
		names:   map[string]string{},
		clashes: map[string]map[string]bool{},
	}
	return s
}

type ScannerStore struct {
	app     *App
	imports map[string]map[string][]string
	names   map[string]string
	clashes map[string]map[string]bool
}

func (s *ScannerStore) Clashes() map[string]map[string]bool {
	return s.clashes
}

// Main is the path of the main package
func (s *ScannerStore) Main() (string, int) {
	if s.names[s.app.Editor.CurrentPackage()] == "main" {
		return s.app.Editor.CurrentPackage(), 0
	}
	// count the main packages
	var count int
	var path string
	for p, n := range s.names {
		if n == "main" {
			count++
			path = p
		}
	}
	if count == 1 {
		return path, 1
	}
	return "", count
}

func (s *ScannerStore) MainPackages() map[string]bool {
	m := map[string]bool{}
	for p, n := range s.names {
		if n == "main" {
			m[p] = true
		}
	}
	return m
}

func (s *ScannerStore) DisplayPath(path string) string {
	parts := strings.Split(path, "/")
	guessed := parts[len(parts)-1]
	name := s.names[path]
	suffix := ""
	if guessed != name && name != "" {
		suffix = " (" + name + ")"
	}
	return path + suffix
}

func (s *ScannerStore) DisplayName(path string) string {
	if s.names[path] != "" {
		return s.names[path]
	}
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}

func (s *ScannerStore) Name(path string) string {
	return s.names[path]
}

func (s *ScannerStore) Names() map[string]string {
	return s.names
}

// Imports returns all the imports from all files in a package
func (s *ScannerStore) Imports(path string) []string {
	var a []string
	for _, f := range s.imports[path] {
		for _, i := range f {
			a = append(a, i)
		}
	}
	return a
}

func (s *ScannerStore) AllImports() map[string]bool {
	m := map[string]bool{}
	for _, imps := range s.imports {
		for _, file := range imps {
			for _, imp := range file {
				m[imp] = true
			}
		}
	}
	return m
}

func (s *ScannerStore) AllImportsOrdered() []string {
	var a []string
	m := map[string]bool{}
	for _, imps := range s.imports {
		for _, file := range imps {
			for _, imp := range file {
				if !m[imp] {
					m[imp] = true
					a = append(a, imp)
				}
			}
		}
	}
	sort.Strings(a)
	return a
}

func (s *ScannerStore) Handle(payload *flux.Payload) bool {
	switch action := payload.Action.(type) {
	case *actions.DeleteFile:
		delete(s.imports[s.app.Editor.CurrentPackage()], action.Name)
		payload.Notify()
	case *actions.RemovePackage:
		payload.Wait(s.app.Source)
		delete(s.imports, action.Path)
		delete(s.names, action.Path)
		s.checkForClash()
		payload.Notify()
	case *actions.AddPackage:
		payload.Wait(s.app.Source)
		if s.checkForClash() {
			payload.Notify()
		}
	case *actions.LoadSource:
		payload.Wait(s.app.Source)

		for path := range action.Source {
			delete(s.imports, path)
			delete(s.names, path)
		}

		var changed bool
		for path, files := range action.Source {
			for name, contents := range files {
				if s.refresh(path, name, contents) {
					changed = true
				}
			}
		}

		if s.checkForClash() {
			changed = true
		}

		if changed {
			payload.Notify()
		}
	case *actions.RequestClose:
		payload.Wait(s.app.Archive)
		if s.checkForClash() {
			payload.Notify()
		}
	case *actions.DragDrop:
		payload.Wait(s.app.Source)
		var changed bool
		for path, files := range action.Changed {
			for name := range files {
				if s.refresh(path, name, s.app.Source.Contents(path, name)) {
					changed = true
				}
			}
		}
		if changed {
			payload.Notify()
		}
	case *actions.UserChangedText:
		payload.Wait(s.app.Source)
		if action.Changed {
			if s.refresh(s.app.Editor.CurrentPackage(), s.app.Editor.CurrentFile(), s.app.Source.Current()) {
				payload.Notify()
			}
		}
	}
	return true
}

func (s *ScannerStore) checkForClash() bool {

	clashes := map[string]map[string]bool{}

	var check func(path string)
	check = func(path string) {
		imports := map[string]bool{}
		if s.app.Source.HasPackage(path) {
			for _, imp := range s.app.Scanner.Imports(path) {
				imports[imp] = true
			}
		} else if ci, ok := s.app.Archive.Cache()[path]; ok {
			for _, imp := range ci.Archive.Imports {
				imports[imp] = true
			}
		}
		for imp := range imports {
			check(imp)
			if !s.app.Source.HasPackage(path) && s.app.Source.HasPackage(imp) {
				if clashes[imp] == nil {
					clashes[imp] = map[string]bool{}
				}
				clashes[imp][path] = true
			}
		}
	}
	for path := range s.MainPackages() {
		check(path)
	}

	hadClashesBefore := len(s.clashes) > 0
	hasClashesAfter := len(clashes) > 0

	s.clashes = clashes

	return hadClashesBefore || hasClashesAfter
}

func (s *ScannerStore) refresh(path, filename, contents string) bool {

	if strings.HasSuffix(filename, "_test.go") {
		return false
	}

	fset := token.NewFileSet()

	// ignore errors
	f, _ := parser.ParseFile(fset, filename, contents, parser.ImportsOnly)

	name := f.Name.Name

	var nameChanged bool
	if s.names[path] != name {
		nameChanged = true
		s.names[path] = name
	}

	var imports []string
	for _, v := range f.Imports {
		// ignore errors
		unquoted, _ := strconv.Unquote(v.Path.Value)
		imports = append(imports, unquoted)
	}
	sort.Strings(imports)

	var importsChanged bool
	if s.imports[path] == nil {
		s.imports[path] = map[string][]string{}
	}
	if s.changed(s.imports[path][filename], imports) {
		importsChanged = true
		s.imports[path][filename] = imports
	}

	return nameChanged || importsChanged
}

func (s *ScannerStore) changed(imports, compare []string) bool {
	if len(compare) != len(imports) {
		return true
	}
	for i := range compare {
		if imports[i] != compare[i] {
			return true
		}
	}
	return false
}
