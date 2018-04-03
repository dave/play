package stores

import (
	"fmt"
	"go/types"

	"bytes"

	"encoding/gob"

	"compress/gzip"

	"errors"

	"net/http"

	"sync"

	"github.com/dave/flux"
	"github.com/dave/jsgo/builderjs"
	"github.com/dave/jsgo/server/messages"
	"github.com/dave/play/actions"
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
	Archive *compiler.Archive
}

func NewArchiveStore(app *App) *ArchiveStore {
	s := &ArchiveStore{
		app:   app,
		cache: map[string]CacheItem{},
	}
	return s
}

func (s *ArchiveStore) Compile(path string) ([]*compiler.Archive, error) {
	done := make(map[string]bool)
	archives := map[string]*compiler.Archive{}
	packages := map[string]*types.Package{}
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
				false,
				archives,
				packages,
			)
			if err != nil {
				return err
			}
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
	return deps, nil
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

func (s *ArchiveStore) Handle(payload *flux.Payload) bool {
	switch a := payload.Action.(type) {
	case *actions.LoadSource:
		payload.Wait(s.app.Scanner)
		if !s.AllFresh() {
			s.app.Dispatch(&actions.UpdateStart{})
		}
	case *actions.UpdateStart:
		path, count := s.app.Scanner.Main()
		if path == "" {
			if count == 0 {
				s.app.Fail(errors.New("project has no main package"))
				return true
			} else {
				s.app.Fail(fmt.Errorf("project has %d main packages - select one and retry", count))
				return true
			}
		}
		s.app.Log("updating")
		s.index = nil
		s.app.Dispatch(&actions.Dial{
			Url:     defaultUrl(),
			Open:    func() flux.ActionInterface { return &actions.UpdateOpen{Main: path} },
			Message: func(m interface{}) flux.ActionInterface { return &actions.UpdateMessage{Message: m} },
			Close:   func() flux.ActionInterface { return &actions.UpdateClose{Run: a.Run, Main: path} },
		})
		payload.Notify()

	case *actions.UpdateOpen:
		hashes := map[string]string{}
		for path, item := range s.Cache() {
			hashes[path] = item.Hash
		}
		message := messages.Update{
			Main:   a.Main,
			Source: s.app.Source.Source(),
			Cache:  hashes,
		}
		s.app.Dispatch(&actions.Send{
			Message: message,
		})
	case *actions.UpdateMessage:
		switch message := a.Message.(type) {
		case messages.Queueing:
			if message.Position > 1 {
				s.app.Logf("queued position %d", message.Position)
			}
		case messages.Downloading:
			if message.Message != "" {
				s.app.Log(message.Message)
			} else if message.Done {
				s.app.Log("building")
			}
		case messages.Archive:
			if message.Standard {
				s.wait.Add(1)
				go func() {
					defer s.wait.Done()
					resp, err := http.Get(fmt.Sprintf("https://%s/%s.%s.a", s.app.PkgHost(), message.Path, message.Hash))
					if err != nil {
						s.app.Fail(err)
						return
					}
					var a compiler.Archive
					if err := gob.NewDecoder(resp.Body).Decode(&a); err != nil {
						s.app.Fail(err)
						return
					}
					s.cache[message.Path] = CacheItem{
						Hash:    message.Hash,
						Archive: &a,
					}
					s.app.Log(a.Name)
				}()
				return true
			} else {
				r, err := gzip.NewReader(bytes.NewBuffer(message.Contents))
				if err != nil {
					s.app.Fail(err)
					return true
				}
				var a compiler.Archive
				if err := gob.NewDecoder(r).Decode(&a); err != nil {
					s.app.Fail(err)
					return true
				}
				s.cache[message.Path] = CacheItem{
					Hash:    message.Hash,
					Archive: &a,
				}
				s.app.Log(a.Name)
			}
		case messages.Index:
			s.index = message
		}
	case *actions.UpdateClose:

		s.wait.Wait()

		if !s.Fresh(a.Main) {
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
