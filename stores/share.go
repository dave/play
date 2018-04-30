package stores

import (
	"fmt"

	"github.com/dave/flux"
	"github.com/dave/jsgo/server/messages"
	"github.com/dave/play/actions"
	"github.com/gopherjs/gopherjs/js"
)

func NewShareStore(app *App) *ShareStore {
	s := &ShareStore{
		app: app,
	}
	return s
}

type ShareStore struct {
	app *App
}

func (s *ShareStore) Handle(payload *flux.Payload) bool {
	switch action := payload.Action.(type) {
	case *actions.ShareStart:
		s.app.Log("sharing")
		s.app.Dispatch(&actions.Dial{
			Url:     defaultUrl(),
			Open:    func() flux.ActionInterface { return &actions.ShareOpen{} },
			Message: func(m interface{}) flux.ActionInterface { return &actions.ShareMessage{Message: m} },
			Close:   func() flux.ActionInterface { return &actions.ShareClose{} },
		})
		payload.Notify()
	case *actions.ShareOpen:
		message := messages.Share{
			Source: s.app.Source.Source(),
			Tags:   s.app.Compile.Tags(),
		}
		s.app.Dispatch(&actions.Send{
			Message: message,
		})
	case *actions.ShareMessage:
		switch message := action.Message.(type) {
		case messages.Storing:
			s.app.Log("storing")
		case messages.ShareComplete:
			js.Global.Get("history").Call("replaceState", js.M{}, "", fmt.Sprintf("/%s", message.Hash))
			s.app.LogHide("shared")
		}
	case *actions.ShareClose:
		// nothing
	}
	return true
}
