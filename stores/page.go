package stores

import (
	"github.com/dave/flux"
	"github.com/dave/play/actions"
	"github.com/dave/play/models"
)

func NewPageStore(app *App) *PageStore {
	s := &PageStore{
		app:      app,
		autoOpen: true,
		modals:   map[models.Modal]bool{},
		minify:   true,
	}
	return s
}

type PageStore struct {
	app         *App
	console     bool
	minify      bool
	autoOpen    bool
	modals      map[models.Modal]bool
	showAllDeps bool // show all dependencies in the load package modal
}

func (s *PageStore) ShowAllDeps() bool {
	return s.showAllDeps
}

func (s *PageStore) ModalOpen(modal models.Modal) bool {
	return s.modals[modal]
}

func (s *PageStore) Console() bool {
	return s.console
}

func (s *PageStore) Minify() bool {
	return s.minify
}

func (s *PageStore) Handle(payload *flux.Payload) bool {
	switch a := payload.Action.(type) {
	case *actions.ModalOpen:
		s.modals[a.Modal] = true
		payload.Notify()
	case *actions.ModalClose:
		s.modals[a.Modal] = false
		payload.Notify()
	case *actions.ConsoleToggleClick:
		s.console = !s.console
		payload.Notify()
	case *actions.MinifyToggleClick:
		s.minify = !s.minify
		payload.Notify()
	case *actions.ConsoleFirstWrite:
		if s.autoOpen {
			s.console = true
			s.autoOpen = false
			payload.Notify()
		}
	case *actions.ShowAllDepsChange:
		s.showAllDeps = a.State
		payload.Notify()
	}
	return true
}
