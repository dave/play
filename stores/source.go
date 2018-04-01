package stores

import (
	"archive/zip"
	"sort"

	"go/format"

	"strings"

	"bytes"
	"io"

	"path/filepath"

	"io/ioutil"

	"errors"

	"fmt"

	"github.com/dave/flux"
	"github.com/dave/play/actions"
	"github.com/dave/saver"
	"github.com/gopherjs/gopherjs/js"
)

func NewSourceStore(app *App) *SourceStore {
	s := &SourceStore{
		app:    app,
		source: map[string]map[string]string{},
	}
	return s
}

type SourceStore struct {
	app *App

	source map[string]map[string]string
}

func (s *SourceStore) Current() string {
	return s.source[s.app.Editor.CurrentPackage()][s.app.Editor.CurrentFile()]
}

func (s *SourceStore) Contents(path, name string) string {
	return s.source[path][name]
}

func (s *SourceStore) Files(path string) map[string]string {
	return s.source[path]
}

func (s *SourceStore) Source() map[string]map[string]string {
	return s.source
}

func (s *SourceStore) HasPackage(path string) bool {
	return s.source[path] != nil
}

func (s *SourceStore) HasFile(path, name string) bool {
	if !s.HasPackage(path) {
		return false
	}
	_, ok := s.source[path][name]
	return ok
}

func (s *SourceStore) Count() int {
	var count int
	for _, pkg := range s.source {
		count += len(pkg)
	}
	return count
}

// SinglePackage returns a random package. Use for when len(s.source) == 1
func (s *SourceStore) SinglePackage() (path string, files map[string]string) {
	for path, files := range s.source {
		return path, files
	}
	return "", nil
}

// SingleFile returns a random file in a random package. Use for when Count() == 1
func (s *SourceStore) SingleFile() (path, name, contents string) {
	for path, files := range s.source {
		for name, contents := range files {
			return path, name, contents
		}
	}
	return "", "", ""
}

func (s *SourceStore) Packages() []string {
	var paths []string
	for p := range s.source {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	return paths
}

func (s *SourceStore) Filenames(path string) []string {
	var f []string
	for k := range s.source[path] {
		f = append(f, k)
	}
	sort.Strings(f)
	return f
}

func (s *SourceStore) Handle(payload *flux.Payload) bool {
	switch a := payload.Action.(type) {
	case *actions.DragEnter:
		s.app.Log("Drop to upload")
	case *actions.DragLeave:
		s.app.Log()
	case *actions.DragDrop:
		s.app.Log()
		packages := map[string]map[string][]byte{}
		if len(a.Files) == 1 && strings.HasSuffix(a.Files[0].Name(), ".zip") {
			b, err := ioutil.ReadAll(a.Files[0].Reader())
			if err != nil {
				s.app.Fail(err)
				return true
			}
			zr, err := zip.NewReader(bytes.NewReader(b), int64(a.Files[0].Len()))
			if err != nil {
				s.app.Fail(err)
				return true
			}
			for _, file := range zr.File {
				path, name := filepath.Split(file.Name)
				path = strings.Trim(path, "/")
				if path == "" {
					path = s.app.Editor.CurrentPackage()
				}
				fr, err := file.Open()
				if err != nil {
					s.app.Fail(err)
					return true
				}
				b, err := ioutil.ReadAll(fr)
				if err != nil {
					fr.Close()
					s.app.Fail(err)
					return true
				}
				fr.Close()
				if packages[path] == nil {
					packages[path] = map[string][]byte{}
				}
				packages[path][name] = b
			}
		} else {
			for _, f := range a.Files {
				path := strings.Trim(f.Dir(), "/")
				if path == "" {
					path = s.app.Editor.CurrentPackage()
				}
				b, err := ioutil.ReadAll(f.Reader())
				if err != nil {
					s.app.Fail(err)
					return true
				}
				if packages[path] == nil {
					packages[path] = map[string][]byte{}
				}
				packages[path][f.Name()] = b
			}
		}

		changed := map[string]map[string]bool{}
		for path, files := range packages {
			for name, contents := range files {
				if !strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, ".jsgo.html") && !strings.HasSuffix(name, ".inc.js") {
					continue
				}
				if s.source[path] == nil {
					s.source[path] = map[string]string{}
				}
				if s.source[path][name] != string(contents) {
					s.source[path][name] = string(contents)

					// track changed files to pass to the scanner
					if changed[path] == nil {
						changed[path] = map[string]bool{}
					}
					changed[path][name] = true
				}
			}
		}

		a.Changed = changed

		payload.Notify()

	case *actions.DownloadClick:
		if s.Count() == 1 {
			_, name, contents := s.SingleFile()
			saver.Save(name, "text/plain", []byte(contents))
			break
		}
		buf := &bytes.Buffer{}
		zw := zip.NewWriter(buf)
		if len(s.source) == 1 {
			_, files := s.SinglePackage()
			for name, contents := range files {
				w, err := zw.Create(name)
				if err != nil {
					s.app.Fail(err)
					return true
				}
				if _, err := io.Copy(w, strings.NewReader(contents)); err != nil {
					s.app.Fail(err)
					return true
				}
			}
		} else {
			for path, files := range s.source {
				for name, contents := range files {
					w, err := zw.Create(filepath.Join(path, name))
					if err != nil {
						s.app.Fail(err)
						return true
					}
					if _, err := io.Copy(w, strings.NewReader(contents)); err != nil {
						s.app.Fail(err)
						return true
					}
				}
			}
		}
		zw.Close()
		saver.Save("src.zip", "application/zip", buf.Bytes())
	case *actions.UserChangedText:
		p := s.app.Editor.CurrentPackage()
		f := s.app.Editor.CurrentFile()
		if p == "" {
			s.app.Fail(errors.New("no package selected"))
			return true
		}
		if f == "" {
			s.app.Fail(errors.New("no file selected"))
			return true
		}
		if s.source[p] == nil {
			s.source[p] = map[string]string{}
		}
		if s.source[p][f] != a.Text {
			s.source[p][f] = a.Text
			a.Changed = true
		}
	case *actions.AddFileClick:
		js.Global.Call("$", "#add-file-modal").Call("modal", "show")
		js.Global.Call("$", "#add-file-input").Call("focus")
		js.Global.Call("$", "#add-file-input").Call("val", "")
		payload.Notify()
	case *actions.AddPackageClick:
		js.Global.Call("$", "#add-package-modal").Call("modal", "show")
		js.Global.Call("$", "#add-package-input").Call("focus")
		js.Global.Call("$", "#add-package-input").Call("val", "")
		payload.Notify()
	case *actions.AddFile:
		js.Global.Call("$", "#add-file-modal").Call("modal", "hide")
		p := s.app.Editor.CurrentPackage()
		if p == "" {
			s.app.Fail(errors.New("no package selected"))
			return true
		}
		if s.source[p] == nil {
			s.source[p] = map[string]string{}
		}
		if s.app.Scanner.Name(p) != "" && strings.HasSuffix(a.Name, ".go") {
			s.source[p][a.Name] = "package " + s.app.Scanner.Name(p) + "\n\n"
		} else {
			s.source[p][a.Name] = ""
		}
		payload.Notify()
	case *actions.AddPackage:
		js.Global.Call("$", "#add-package-modal").Call("modal", "hide")
		if s.source[a.Path] == nil {
			s.source[a.Path] = map[string]string{}
		}
		payload.Notify()
	case *actions.DeleteFileClick:
		js.Global.Call("$", "#delete-file-modal").Call("modal", "show")
		payload.Notify()
	case *actions.RemovePackageClick:
		js.Global.Call("$", "#remove-package-modal").Call("modal", "show")
		payload.Notify()
	case *actions.DeleteFile:
		js.Global.Call("$", "#delete-file-modal").Call("modal", "hide")
		p := s.app.Editor.CurrentPackage()
		if p == "" {
			s.app.Fail(errors.New("no package selected"))
			return true
		}
		if !s.HasPackage(p) {
			s.app.Fail(fmt.Errorf("package %s not found", p))
			return true
		}
		if !s.HasFile(p, a.Name) {
			s.app.Fail(fmt.Errorf("%s not found", a.Name))
			return true
		}
		delete(s.source[p], a.Name)
		payload.Notify()
	case *actions.RemovePackage:
		js.Global.Call("$", "#remove-package-modal").Call("modal", "hide")
		if !s.HasPackage(a.Path) {
			s.app.Fail(fmt.Errorf("%s not found", a.Path))
			return true
		}
		delete(s.source, a.Path)
		payload.Notify()
	case *actions.LoadSource:
		s.source = a.Source
		payload.Notify()
	case *actions.FormatCode:
		p := s.app.Editor.CurrentPackage()
		f := s.app.Editor.CurrentFile()
		if strings.HasSuffix(f, ".go") {
			b, err := format.Source([]byte(s.Contents(p, f)))
			if err != nil {
				s.app.Fail(err)
				return true
			}
			s.source[p][f] = string(b)
			payload.Notify()
		}
		if a.Then != nil {
			s.app.Dispatch(a.Then)
		}
	}
	return true
}
