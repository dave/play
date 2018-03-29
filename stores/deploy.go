package stores

import (
	"github.com/dave/flux"
	"github.com/dave/jsgo/server/messages"
	"github.com/dave/play/actions"
)

func NewDeployStore(app *App) *DeployStore {
	s := &DeployStore{
		app: app,
	}
	return s
}

type DeployStore struct {
	app *App
}

func (s *DeployStore) Handle(payload *flux.Payload) bool {
	switch action := payload.Action.(type) {
	case *actions.DeployStart:
		s.app.Log("deploying")
		s.app.Dispatch(&actions.Dial{
			Url:     defaultUrl(),
			Open:    func() flux.ActionInterface { return &actions.DeployOpen{} },
			Message: func(m interface{}) flux.ActionInterface { return &actions.DeployMessage{Message: m} },
			Close:   func() flux.ActionInterface { return &actions.DeployClose{} },
		})
		payload.Notify()
	case *actions.DeployOpen:
		message := messages.Share{
			Source: map[string]map[string]string{
				"main": s.app.Editor.Files(),
			},
		}
		s.app.Dispatch(&actions.Send{
			Message: message,
		})
	case *actions.DeployMessage:
		switch message := action.Message.(type) {
		case messages.Storing:
			s.app.Log("storing")
		case messages.Complete:
			_ = message // TODO
			s.app.Log("deployed")
		}
	case *actions.DeployClose:
		s.app.Log()
	}
	return true
}
