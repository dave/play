package stores

import (
	"fmt"
	"net/http"
	"strings"

	"encoding/json"

	"regexp"

	"github.com/dave/flux"
	"github.com/dave/jsgo/config"
	"github.com/dave/jsgo/server/play/messages"
	"github.com/dave/locstor"
	"github.com/dave/play/actions"
	"github.com/dave/play/models"
	"honnef.co/go/js/dom"
)

type LocalStore struct {
	app *App

	local       *locstor.DataStore
	initialized bool
}

func NewLocalStore(app *App) *LocalStore {
	s := &LocalStore{
		app:   app,
		local: locstor.NewDataStore(locstor.JSONEncoding),
	}
	return s
}

func (s *LocalStore) Initialized() bool {
	return s.initialized
}

func (s *LocalStore) Handle(payload *flux.Payload) bool {
	switch action := payload.Action.(type) {
	case *actions.Load:
		var sizes []float64
		found, err := s.local.Find("split-sizes", &sizes)
		if err != nil {
			s.app.Fail(err)
			return true
		}
		if !found {
			sizes = defaultSizes
		}
		s.app.Dispatch(&actions.ChangeSplit{Sizes: sizes})

		location := strings.Trim(dom.GetWindow().Location().Pathname, "/")

		var seenHelp bool
		if _, err := s.local.Find("seen-help", &seenHelp); err != nil {
			s.app.Fail(err)
			return true
		}
		if !seenHelp {
			s.app.Dispatch(&actions.ModalOpen{Modal: models.HelpModal})
			if err := s.local.Save("seen-help", true); err != nil {
				s.app.Fail(err)
				return true
			}
		}

		// No page path -> load files from local storage or use default files
		if location == "" {
			var currentPackage, currentFile string
			var buildTags []string
			var source map[string]map[string]string
			found, err = s.local.Find("source", &source)
			if err != nil {
				s.app.Fail(err)
				return true
			}
			if !found {
				// old format for storing files
				var files map[string]string
				found, err = s.local.Find("files", &files)
				if err != nil {
					s.app.Fail(err)
					return true
				}
				if found {
					source = map[string]map[string]string{"main": files}
				}
			}
			if found {
				// if we found files in local storage, also load the current file and package
				if _, err := s.local.Find("current-file", &currentFile); err != nil {
					s.app.Fail(err)
					return true
				}
				if _, err := s.local.Find("current-package", &currentPackage); err != nil {
					s.app.Fail(err)
					return true
				}
				if _, err := s.local.Find("build-tags", &buildTags); err != nil {
					s.app.Fail(err)
					return true
				}
			} else {
				// if we didn't find source in the local storage, add the default file
				source = map[string]map[string]string{"main": {"main.go": defaultFile}}
				currentFile = "main.go"
				currentPackage = "main"
			}
			s.app.Dispatch(&actions.LoadSource{
				Source:         source,
				CurrentFile:    currentFile,
				CurrentPackage: currentPackage,
				Tags:           buildTags,
				Update:         true,
			})
			break
		}

		// Hash in page path -> load files from src.jsgo.io json blob
		if shaRegex.MatchString(location) {
			resp, err := http.Get(fmt.Sprintf("%s://%s/%s.json", config.Protocol[config.Src], config.Host[config.Src], location))
			if err != nil {
				s.app.Fail(err)
				return true
			}
			if resp.StatusCode != 200 {
				s.app.Fail(fmt.Errorf("error %d loading source", resp.StatusCode))
				return true
			}
			var m messages.Share
			if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
				s.app.Fail(err)
				return true
			}
			s.app.Dispatch(&actions.LoadSource{
				Source: m.Source,
				Tags:   m.Tags,
				Update: true,
			})
			break
		}

		// Package path in page path -> open websocket and load files
		s.app.Dispatch(&actions.RequestStart{Type: models.InitialiseRequest, Path: location})

	case *actions.UserChangedSplit:
		if err := s.saveSplitSizes(action.Sizes); err != nil {
			s.app.Fail(err)
			return true
		}
	case *actions.UserChangedText, *actions.FormatCode:
		payload.Wait(s.app.Editor)
		if err := s.saveSource(); err != nil {
			s.app.Fail(err)
			return true
		}
	case *actions.BuildTags:
		payload.Wait(s.app.Compile)
		if err := s.saveSource(); err != nil {
			s.app.Fail(err)
			return true
		}
	case *actions.UserChangedFile:
		payload.Wait(s.app.Editor)
		if err := s.saveSource(); err != nil {
			s.app.Fail(err)
			return true
		}
	case *actions.UserChangedPackage:
		payload.Wait(s.app.Editor)
		if err := s.saveSource(); err != nil {
			s.app.Fail(err)
			return true
		}
	case *actions.AddFile, *actions.DeleteFile:
		payload.Wait(s.app.Source)
		if err := s.saveSource(); err != nil {
			s.app.Fail(err)
			return true
		}
	case *actions.AddPackage, *actions.RemovePackage, *actions.DragDrop:
		payload.Wait(s.app.Editor)
		if err := s.saveSource(); err != nil {
			s.app.Fail(err)
			return true
		}
	case *actions.LoadSource:
		if action.Save {
			payload.Wait(s.app.Editor)
			if err := s.saveSource(); err != nil {
				s.app.Fail(err)
				return true
			}
		}
	}
	return true
}

func (s *LocalStore) saveSource() error {
	s.local.Delete("files") // delete old format file storage location
	if err := s.local.Save("source", s.app.Source.Source()); err != nil {
		return err
	}
	if err := s.local.Save("current-package", s.app.Editor.CurrentPackage()); err != nil {
		return err
	}
	if err := s.local.Save("current-file", s.app.Editor.CurrentFile()); err != nil {
		return err
	}
	if err := s.local.Save("build-tags", s.app.Compile.Tags()); err != nil {
		return err
	}
	return nil
}

func (s *LocalStore) saveSplitSizes(sizes []float64) error {
	return s.local.Save("split-sizes", sizes)
}

var (
	defaultSizes = []float64{50, 50}
	defaultFile  = `package main

import (
	"fmt"
	"honnef.co/go/js/dom"
)

func main() {
	body := dom.GetWindow().Document().GetElementsByTagName("body")[0]
	body.SetInnerHTML("Hello, HTML!")
	fmt.Println("Hello, console!")
}`
)

var shaRegex = regexp.MustCompile("^[0-9a-f]{40}$")
