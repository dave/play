package stores

import (
	"github.com/dave/flux"
	"github.com/dave/play/actions"
)

func NewEditorStore(app *App) *EditorStore {
	s := &EditorStore{
		app:          app,
		currentFiles: map[string]string{},
	}
	return s
}

type EditorStore struct {
	app *App

	sizes          []float64
	currentPackage string
	currentFiles   map[string]string // tracks the currently selected file in each package
	loaded         bool
}

func (s *EditorStore) Loaded() bool {
	return s.loaded
}

func (s *EditorStore) Sizes() []float64 {
	return s.sizes
}

func (s *EditorStore) CurrentPackage() string {
	return s.currentPackage
}

func (s *EditorStore) CurrentFile() string {
	return s.currentFiles[s.currentPackage]
}

func (s *EditorStore) defaultPackage() string {
	if len(s.app.Source.Packages()) == 0 {
		return ""
	}
	var path string

	// default to the first main package
	for p, n := range s.app.Scanner.Names() {
		if n == "main" {
			path = p
			break
		}
	}

	if path == "" {
		// if no main package, choose first package ordered by name
		path = s.app.Source.Packages()[0]
	}

	return path
}

func (s *EditorStore) defaultFile(path string) string {
	if len(s.app.Source.Files(path)) == 0 {
		return ""
	}
	if _, ok := s.app.Source.Files(path)["README.md"]; ok {
		return "README.md"
	}
	if _, ok := s.app.Source.Files(path)["readme.md"]; ok {
		return "readme.md"
	}
	if _, ok := s.app.Source.Files(path)["main.go"]; ok {
		return "main.go"
	}
	return s.app.Source.Filenames(path)[0]
}

func (s *EditorStore) Handle(payload *flux.Payload) bool {
	switch a := payload.Action.(type) {
	case *actions.DragDrop:
		payload.Wait(s.app.Source)
		s.currentPackage = s.defaultPackage()
		s.currentFiles[s.currentPackage] = s.defaultFile(s.currentPackage)
		payload.Notify()
	case *actions.LoadSource:
		payload.Wait(s.app.Scanner)

		defer func() {
			s.loaded = true
		}()

		var switchPackage string
		for path := range a.Source {
			switchPackage = path
			if s.app.Scanner.Name(path) == "main" {
				// break if we find a main package
				break
			}
		}

		// Switch to the right package.
		if a.CurrentPackage != "" && s.app.Source.HasPackage(a.CurrentPackage) {
			s.currentPackage = a.CurrentPackage
		} else {
			s.currentPackage = switchPackage
		}

		// Switch to the right file.
		if a.CurrentFile != "" && s.app.Source.HasFile(s.currentPackage, a.CurrentFile) {
			s.currentFiles[s.currentPackage] = a.CurrentFile
		} else {
			s.currentFiles[s.currentPackage] = s.defaultFile(s.currentPackage)
		}

		payload.Notify()

	case *actions.ChangeSplit:
		s.sizes = a.Sizes
		payload.Notify()
	case *actions.UserChangedSplit:
		s.sizes = a.Sizes
	case *actions.UserChangedFile:
		s.currentFiles[s.currentPackage] = a.Name
		payload.Notify()
	case *actions.ChangeFile:
		s.currentPackage = a.Path
		s.currentFiles[a.Path] = a.Name
		payload.Notify()
	case *actions.UserChangedPackage:
		s.currentPackage = a.Path
		if s.currentFiles[a.Path] == "" {
			s.currentFiles[a.Path] = s.defaultFile(a.Path)
		}
		payload.Notify()
	case *actions.AddFile:
		payload.Wait(s.app.Source)
		s.currentFiles[s.currentPackage] = a.Name
		payload.Notify()
	case *actions.DeleteFile:
		payload.Wait(s.app.Source)
		if s.CurrentFile() == a.Name {
			s.currentFiles[s.currentPackage] = s.defaultFile(s.currentPackage)
			payload.Notify()
		}
	case *actions.AddPackage:
		payload.Wait(s.app.Source)
		s.currentPackage = a.Path
		payload.Notify()
	case *actions.RemovePackage:
		payload.Wait(s.app.Scanner)
		if s.currentPackage == a.Path {
			s.currentPackage = s.defaultPackage()
			if s.currentFiles[s.currentPackage] == "" {
				s.currentFiles[s.currentPackage] = s.defaultFile(s.currentPackage)
			}
			payload.Notify()
		}
	}
	return true
}
