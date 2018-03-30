package stores

import (
	"archive/zip"
	"sort"

	"errors"

	"go/format"

	"strings"

	"bytes"
	"io"

	"path/filepath"

	"io/ioutil"

	"github.com/dave/flux"
	"github.com/dave/play/actions"
	"github.com/dave/saver"
	"github.com/gopherjs/gopherjs/js"
)

func NewEditorStore(app *App) *EditorStore {
	s := &EditorStore{
		app:   app,
		files: map[string]string{},
	}
	return s
}

type EditorStore struct {
	app *App

	sizes   []float64
	files   map[string]string
	current string
}

func (s *EditorStore) Sizes() []float64 {
	return s.sizes
}

func (s *EditorStore) Text() string {
	return s.files[s.current]
}

func (s *EditorStore) Current() string {
	return s.current
}

func (s *EditorStore) Files() map[string]string {
	f := map[string]string{}
	for k, v := range s.files {
		f[k] = v
	}
	return f
}

func (s *EditorStore) Filenames() []string {
	var f []string
	for k := range s.files {
		f = append(f, k)
	}
	sort.Strings(f)
	return f
}

func (s *EditorStore) Handle(payload *flux.Payload) bool {
	switch a := payload.Action.(type) {
	case *actions.DragEnter:
		s.app.Log("Drop to upload")
	case *actions.DragLeave:
		s.app.Log()
	case *actions.DragDrop:
		s.app.Log()
		// TODO: support multiple packages
		files := map[string][]byte{}
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
				_, name := filepath.Split(file.Name)
				if files[name] != nil {
					// two files might have the same name in different dirs
					continue
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
				files[name] = b
			}
		} else {
			for _, f := range a.Files {
				if files[f.Name()] != nil {
					// two files might have the same name in different dirs
					continue
				}
				b, err := ioutil.ReadAll(f.Reader())
				if err != nil {
					s.app.Fail(err)
					return true
				}
				files[f.Name()] = b
			}
		}

		for name, contents := range files {
			if !strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, ".jsgo.html") && !strings.HasSuffix(name, ".inc.js") {
				continue
			}
			s.files[name] = string(contents)
			if len(files) == 1 {
				s.current = name
			}
		}
		payload.Notify()

	case *actions.DownloadClick:
		if len(s.files) == 1 {
			var name, contents string
			for n, c := range s.files {
				name = n
				contents = c
			}
			saver.Save(name, "text/plain", []byte(contents))
			break
		}

		buf := &bytes.Buffer{}
		zw := zip.NewWriter(buf)
		for name, contents := range s.files {
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
		zw.Close()
		saver.Save("source.zip", "application/zip", buf.Bytes())
	case *actions.ChangeSplit:
		s.sizes = a.Sizes
		payload.Notify()
	case *actions.ChangeText:
		s.files[s.current] = a.Text
		payload.Notify()
	case *actions.UserChangedSplit:
		s.sizes = a.Sizes
	case *actions.UserChangedText:
		s.clearState()
		s.files[s.current] = a.Text
	case *actions.UserChangedFile:
		s.current = a.Name
		payload.Notify()
	case *actions.AddFileClick:
		js.Global.Call("$", "#add-file-modal").Call("modal", "show")
		js.Global.Call("$", "#add-file-input").Call("focus")
		payload.Notify()
	case *actions.DeleteFileClick:
		js.Global.Call("$", "#delete-file-modal").Call("modal", "show")
		js.Global.Call("$", "#delete-file-input").Call("focus")
		payload.Notify()
	case *actions.AddFile:
		js.Global.Call("$", "#add-file-modal").Call("modal", "hide")
		if s.app.Scanner.Name() != "" && strings.HasSuffix(a.Name, ".go") {
			s.files[a.Name] = "package " + s.app.Scanner.Name() + "\n\n"
		} else {
			s.files[a.Name] = ""
		}
		s.clearState()
		s.current = a.Name
		payload.Notify()
	case *actions.DeleteFile:
		js.Global.Call("$", "#delete-file-modal").Call("modal", "hide")
		if len(s.files) == 1 {
			s.app.Fail(errors.New("can't delete last file"))
			return true
		}
		delete(s.files, a.Name)
		if s.current == a.Name {
			s.current = s.Filenames()[0]
		}
		s.clearState()
		payload.Notify()
	case *actions.LoadFiles:
		s.files = a.Files
		var found bool
		current := a.Current
		// if no current file specified, default to "main.go". If it doesn't exist, use the first file.
		if current == "" {
			current = "main.go"
		}
		for name := range s.files {
			if current == name {
				found = true
				s.current = current
				break
			}
		}
		if !found && len(s.files) > 0 {
			s.current = s.Filenames()[0]
		}
		s.app.Dispatch(&actions.ChangeText{
			Text: s.files[s.current],
		})
		payload.Notify()
	case *actions.FormatCode:
		if strings.HasSuffix(s.current, ".go") {
			b, err := format.Source([]byte(s.files[s.current]))
			if err != nil {
				s.app.Fail(err)
				return true
			}
			s.files[s.current] = string(b)
			payload.Notify()
		}
		if a.Then != nil {
			s.app.Dispatch(a.Then)
		}
	}
	return true
}

func (s *EditorStore) clearState() {
	js.Global.Get("history").Call("replaceState", js.M{}, "", "/")
}
