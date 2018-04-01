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
	}
	return s
}

type ScannerStore struct {
	app     *App
	imports map[string]map[string][]string
	names   map[string]string
	main    string
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

func (s *ScannerStore) Handle(payload *flux.Payload) bool {
	switch action := payload.Action.(type) {
	case *actions.DeleteFile:
		delete(s.imports[s.app.Editor.CurrentPackage()], action.Name)
		payload.Notify()
	case *actions.RemovePackage:
		delete(s.imports, action.Path)
		delete(s.names, action.Path)
		if s.main == action.Path {
			// If we delete the main package, try to find another main package
			s.main = ""
			for path, name := range s.names {
				if name == "main" {
					s.main = path
					break
				}
			}
		}
		payload.Notify()
	case *actions.LoadSource:
		payload.Wait(s.app.Source)

		// source is replaced, so clear all imports
		s.imports = map[string]map[string][]string{}
		s.names = map[string]string{}
		s.main = ""

		var changed bool
		for path, files := range s.app.Source.Source() {
			for name, contents := range files {
				if s.refresh(path, name, contents) {
					changed = true
				}
			}
		}
		if changed {
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

func (s *ScannerStore) refresh(path, filename, contents string) bool {
	fset := token.NewFileSet()

	// ignore errors
	f, _ := parser.ParseFile(fset, filename, contents, parser.ImportsOnly)

	name := strings.TrimSuffix(f.Name.Name, "_test")
	var nameChanged bool
	if s.names[path] != name {
		nameChanged = true
		s.names[path] = name
	}

	var mainChanged bool
	if name == "main" && s.main != path {
		mainChanged = true
		s.main = path
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

	return nameChanged || importsChanged || mainChanged
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
