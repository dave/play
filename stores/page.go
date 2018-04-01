package stores

import (
	"github.com/dave/flux"
	"github.com/dave/play/actions"
)

func NewPageStore(app *App) *PageStore {
	s := &PageStore{
		app:      app,
		autoOpen: true,
	}
	return s
}

type PageStore struct {
	app      *App
	console  bool
	autoOpen bool
}

func (s *PageStore) Console() bool {
	return s.console
}
func (s *PageStore) Handle(payload *flux.Payload) bool {
	switch payload.Action.(type) {
	case *actions.ConsoleToggleClick:
		s.console = !s.console
		payload.Notify()
	case *actions.ConsoleFirstWrite:
		if s.autoOpen {
			s.console = true
			s.autoOpen = false
			payload.Notify()
		}
	}
	return true
}
