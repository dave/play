package stores

import (
	"github.com/dave/flux"
	"github.com/dave/play/actions"
	"github.com/gopherjs/gopherjs/js"
)

func NewHistoryStore(app *App) *HistoryStore {
	s := &HistoryStore{
		app: app,
	}
	return s
}

type HistoryStore struct {
	app *App
}

func (s *HistoryStore) Handle(payload *flux.Payload) bool {
	switch a := payload.Action.(type) {
	case *actions.UserChangedText,
		*actions.AddFile,
		*actions.DeleteFile,
		*actions.AddPackage,
		*actions.RemovePackage,
		*actions.DragDrop,
		*actions.BuildTags:
		js.Global.Get("history").Call("replaceState", js.M{}, "", "/")
	case *actions.LoadSource:
		if a.Save {
			js.Global.Get("history").Call("replaceState", js.M{}, "", "/")
		}
	}
	return true
}
