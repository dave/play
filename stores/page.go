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
	}
	return s
}

type PageStore struct {
	app      *App
	console  bool
	autoOpen bool
	modals   map[models.Modal]bool
}

func (s *PageStore) ModalOpen(modal models.Modal) bool {
	return s.modals[modal]
}

func (s *PageStore) Console() bool {
	return s.console
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
	case *actions.ConsoleFirstWrite:
		if s.autoOpen {
			s.console = true
			s.autoOpen = false
			payload.Notify()
		}
	}
	return true
}
