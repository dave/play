package stores

import (
	"errors"
	"fmt"

	"github.com/dave/flux"
	"github.com/dave/jsgo/server/messages"
	"github.com/dave/play/actions"
	"github.com/dave/play/models"
)

func NewDeployStore(app *App) *DeployStore {
	s := &DeployStore{
		app: app,
	}
	return s
}

type DeployStore struct {
	app                           *App
	mainHash, indexHash, mainPath string
}

func (s *DeployStore) LoaderJs() string {
	return fmt.Sprintf("%s://%s/%s.%s.js", s.app.Protocol(), s.app.PkgHost(), s.mainPath, s.mainHash)
}

func (s *DeployStore) Index() string {
	return fmt.Sprintf("%s://%s/%s", s.app.Protocol(), s.app.IndexHost(), s.indexHash)
}

func (s *DeployStore) Handle(payload *flux.Payload) bool {
	switch action := payload.Action.(type) {
	case *actions.DeployStart:
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
		s.app.Log("deploying")
		s.mainHash = ""
		s.indexHash = ""
		s.mainPath = path
		s.app.Dispatch(&actions.Dial{
			Url:     defaultUrl(),
			Open:    func() flux.ActionInterface { return &actions.DeployOpen{} },
			Message: func(m interface{}) flux.ActionInterface { return &actions.DeployMessage{Message: m} },
			Close:   func() flux.ActionInterface { return &actions.DeployClose{} },
		})
		payload.Notify()
	case *actions.DeployOpen:
		message := messages.Deploy{
			Main:   s.mainPath,
			Source: s.app.Source.Source(),
		}
		s.app.Dispatch(&actions.Send{
			Message: message,
		})
	case *actions.DeployMessage:
		switch message := action.Message.(type) {
		case messages.Downloading:
			if message.Message != "" {
				s.app.Log(message.Message)
			}
		case messages.Compiling:
			if message.Message != "" {
				s.app.Log(message.Message)
			}
		case messages.Storing:
			s.app.Log("storing")
		case messages.DeployComplete:
			s.mainHash = message.Main
			s.indexHash = message.Index
			s.app.Dispatch(&actions.ModalOpen{Modal: models.DeployDoneModal})
			s.app.Log("deployed")
			payload.Notify()
		}
	case *actions.DeployClose:
		s.app.Log()
	}
	return true
}
