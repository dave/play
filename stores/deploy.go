package stores

import (
	"fmt"

	"github.com/dave/flux"
	"github.com/dave/jsgo/server/messages"
	"github.com/dave/play/actions"
	"github.com/gopherjs/gopherjs/js"
)

func NewDeployStore(app *App) *DeployStore {
	s := &DeployStore{
		app: app,
	}
	return s
}

type DeployStore struct {
	app                 *App
	mainHash, indexHash string
}

func (s *DeployStore) LoaderJs() string {
	return fmt.Sprintf("%s://%s/%s.%s.js", s.app.Protocol(), s.app.PkgHost(), "main", s.mainHash)
}

func (s *DeployStore) Index() string {
	return fmt.Sprintf("%s://%s/%s", s.app.Protocol(), s.app.IndexHost(), s.indexHash)
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
		message := messages.Deploy{
			Source: map[string]map[string]string{
				"main": s.app.Editor.Files(),
			},
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
			js.Global.Call("$", "#deploy-done-modal").Call("modal", "show")
			s.app.Log("deployed")
			payload.Notify()
		}
	case *actions.DeployClose:
		s.app.Log()
	}
	return true
}
