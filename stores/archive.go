package stores

import (
	"context"
	"fmt"
	"go/types"

	"encoding/gob"

	"errors"

	"net/http"

	"sync"

	"io/ioutil"

	"github.com/dave/flux"
	"github.com/dave/jsgo/config"
	"github.com/dave/jsgo/server/messages"
	"github.com/dave/play/actions"
	"github.com/dave/play/models"
	"github.com/dave/play/stores/builderjs"
	"github.com/gopherjs/gopherjs/compiler"
)

type ArchiveStore struct {
	app *App

	// cache (path -> item) of archives
	cache map[string]CacheItem

	// index (path -> item) of the previously received update
	index messages.Index

	wait sync.WaitGroup
}

type CacheItem struct {
	Hash    string
	Archive *compiler.Archive // This archive is stripped of JS
	Js      []byte
}

func NewArchiveStore(app *App) *ArchiveStore {
	s := &ArchiveStore{
		app:   app,
		cache: map[string]CacheItem{},
	}
	return s
}

type Dep struct {
	Path string
	Js   []byte
}

func (s *ArchiveStore) Compile(path string) ([]Dep, error) {
	done := make(map[string]bool)
	archives := map[string]*compiler.Archive{}
	packages := map[string]*types.Package{}
	jsdeps := []Dep{
		// Always start with the prelude
		{Path: "prelude", Js: s.cache["prelude"].Js},
	}
	var deps []*compiler.Archive
	var compile func(path string) error
	compile = func(path string) error {
		if done[path] {
			return nil
		}
		if s.app.Source.HasPackage(path) {
			for _, imp := range s.app.Scanner.Imports(path) {
				if err := compile(imp); err != nil {
					return err
				}
			}
			archive, err := builderjs.BuildPackage(
				path,
				s.app.Source.Source(),
				deps,
				s.app.Page.Minify(),
				archives,
				packages,
			)
			if err != nil {
				return err
			}
			js, _, err := builderjs.GetPackageCode(context.Background(), archive, false, true)
			if err != nil {
				return err
			}
			jsdeps = append(jsdeps, Dep{path, js})
			deps = append(deps, archive)
			done[path] = true
			return nil
		}
		item, ok := s.cache[path]
		if !ok {
			return fmt.Errorf("%s not found", path)
		}
		for _, imp := range item.Archive.Imports {
			if err := compile(imp); err != nil {
				return err
			}
		}
		jsdeps = append(jsdeps, Dep{path, item.Js})
		deps = append(deps, item.Archive)
		done[path] = true
		return nil
	}
	if err := compile("runtime"); err != nil {
		return nil, err
	}
	if err := compile(path); err != nil {
		return nil, err
	}
	return jsdeps, nil
}

func (s *ArchiveStore) AllFresh() bool {
	for path := range s.app.Scanner.MainPackages() {
		if !s.Fresh(path) {
			return false
		}
	}
	return true
}

// Fresh is true if current cache matches the previously downloaded archives
func (s *ArchiveStore) Fresh(mainPath string) bool {
	// if index is nil, either the page has just loaded or we're in the middle of an update
	if s.index == nil {
		return false
	}

	// first check that all indexed packages are in the cache at the right versions. This would fail
	// if there was an error while downloading one of the archive files.
	for path, item := range s.index {
		cached, ok := s.cache[path]
		if !ok {
			return false
		}
		if cached.Hash != item.Hash {
			return false
		}
	}

	// then check that all the imports in all packages are found in the index, or in the source
	for _, path := range s.app.Scanner.Imports(mainPath) {
		_, inIndex := s.index[path]
		_, inSource := s.app.Source.Source()[path]
		if !inIndex && !inSource {
			return false
		}
	}

	return true
}

func (s *ArchiveStore) Cache() map[string]CacheItem {
	return s.cache
}

func (s *ArchiveStore) CacheStrings() map[string]string {
	hashes := map[string]string{}
	for path, item := range s.app.Archive.Cache() {
		hashes[path] = item.Hash
	}
	return hashes
}

func (s *ArchiveStore) Handle(payload *flux.Payload) bool {
	switch a := payload.Action.(type) {
	case *actions.MinifyToggleClick:
		payload.Wait(s.app.Page)
		s.index = nil
		s.app.Dispatch(&actions.RequestStart{Type: models.UpdateRequest, Run: false})
	case *actions.LoadSource:
		payload.Wait(s.app.Scanner)
		if a.Update && !s.AllFresh() {
			s.app.Dispatch(&actions.RequestStart{Type: models.UpdateRequest, Run: false})
		}
	case *actions.RequestMessage:
		switch message := a.Message.(type) {
		case messages.Archive:
			s.wait.Add(1)
			go func() {
				defer s.wait.Done()
				c := CacheItem{
					Hash: message.Hash,
				}
				var getwait sync.WaitGroup
				getwait.Add(2)
				go func() {
					defer getwait.Done()
					if message.Path == "prelude" {
						// prelude doesn't have an archive file
						return
					}
					resp, err := http.Get(fmt.Sprintf("%s://%s/%s.%s.ax", config.Protocol, config.PkgHost, message.Path, message.Hash))
					if err != nil {
						s.app.Fail(err)
						return
					}
					var a compiler.Archive
					if err := gob.NewDecoder(resp.Body).Decode(&a); err != nil {
						s.app.Fail(err)
						return
					}
					c.Archive = &a
				}()
				go func() {
					defer getwait.Done()
					resp, err := http.Get(fmt.Sprintf("%s://%s/%s.%s.js", config.Protocol, config.PkgHost, message.Path, message.Hash))
					if err != nil {
						s.app.Fail(err)
						return
					}
					js, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						s.app.Fail(err)
						return
					}
					c.Js = js
				}()
				getwait.Wait()
				s.cache[message.Path] = c
				if message.Path == "prelude" {
					// prelude doesn't have an archive file
					s.app.Log("prelude")
				} else {
					s.app.Log(c.Archive.Name)
				}
			}()
			return true
		case messages.Index:
			s.index = message
		}
	case *actions.RequestClose:

		if a.Type == models.GetRequest {
			// get request doesn't do an update - just gets files
			return true
		}

		s.wait.Wait()

		if !s.AllFresh() {
			s.app.Fail(errors.New("websocket closed but archives not updated"))
			return true
		}

		if a.Run {
			s.app.Dispatch(&actions.CompileStart{})
		} else {
			var downloaded, unchanged int
			for _, v := range s.index {
				if v.Unchanged {
					unchanged++
				} else {
					downloaded++
				}
			}
			if downloaded == 0 && unchanged == 0 {
				s.app.Log()
			} else if downloaded > 0 && unchanged > 0 {
				s.app.LogHidef("%d downloaded, %d unchanged", downloaded, unchanged)
			} else if downloaded > 0 {
				s.app.LogHidef("%d downloaded", downloaded)
			} else if unchanged > 0 {
				s.app.LogHidef("%d unchanged", unchanged)
			}
		}
		payload.Notify()
	}

	return true
}
