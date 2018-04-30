package stores

import (
	"github.com/dave/flux"
	"github.com/dave/jsgo/server/messages"
	"github.com/dave/play/actions"
	"github.com/dave/play/models"
)

func NewRequestStore(app *App) *RequestStore {
	s := &RequestStore{
		app: app,
	}
	return s
}

type RequestStore struct {
	app *App
}

func (s *RequestStore) Handle(payload *flux.Payload) bool {
	switch action := payload.Action.(type) {
	case *actions.RequestStart:
		s.app.Log("downloading")
		s.app.Dispatch(&actions.Dial{
			Url:  defaultUrl(),
			Open: func() flux.ActionInterface { return &actions.RequestOpen{RequestStart: action} },
			Message: func(m interface{}) flux.ActionInterface {
				return &actions.RequestMessage{RequestStart: action, Message: m}
			},
			Close: func() flux.ActionInterface { return &actions.RequestClose{RequestStart: action} },
		})
		payload.Notify()
	case *actions.RequestOpen:
		var message messages.Message
		switch action.Type {
		case models.GetRequest:
			message = messages.Get{
				Path: action.Path,
			}
		case models.UpdateRequest:
			message = messages.Update{
				Source: s.app.Source.Source(),
				Cache:  s.app.Archive.CacheStrings(),
				Minify: s.app.Page.Minify(),
				Tags:   s.app.Compile.Tags(),
			}
		case models.InitialiseRequest:
			message = messages.Initialise{
				Path:   action.Path,
				Minify: s.app.Page.Minify(),
			}
		}
		s.app.Dispatch(&actions.Send{
			Message: message,
		})
	case *actions.RequestMessage:
		switch message := action.Message.(type) {
		case messages.Queueing:
			if message.Position > 1 {
				s.app.Logf("queued position %d", message.Position)
			}
		case messages.Downloading:
			if len(message.Message) > 0 {
				s.app.Log(message.Message)
			}
		case messages.GetComplete:
			s.app.Dispatch(&actions.LoadSource{
				Source: message.Source,
				Save:   action.Type == models.GetRequest,
				Update: action.Type != models.InitialiseRequest, // never update after initialise (will be updated by initialise)
			})
			var count int
			for _, files := range message.Source {
				count += len(files)
			}
			if count == 1 {
				s.app.LogHide("got 1 file")
			} else {
				s.app.LogHidef("got %d files", count)
			}
		}
	case *actions.RequestClose:
		// nothing
	}
	return true
}
