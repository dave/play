package views

import (
	"time"

	"github.com/dave/play/actions"
	"github.com/dave/play/stores"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
	"github.com/tulir/gopher-ace"
	"honnef.co/go/js/dom"
)

type Editor struct {
	vecty.Core
	app *stores.App

	editor ace.Editor
}

func NewEditor(app *stores.App) *Editor {
	v := &Editor{
		app: app,
	}
	return v
}

func (v *Editor) Mount() {
	v.app.Watch(v, func(done chan struct{}) {
		defer close(done)
		if v.app.Source.Current() != v.editor.GetValue() {
			// only update the editor if the text is changed
			v.editor.SetValue(v.app.Source.Current())
			v.editor.ClearSelection()
			v.editor.MoveCursorTo(0, 0)
		}
	})

	v.editor = ace.Edit("editor")
	v.editor.SetOptions(map[string]interface{}{
		"mode": "ace/mode/golang",
	})

	dom.GetWindow().AddEventListener("resize", false, func(event dom.Event) {
		v.Resize()
	})

	v.editor.Get("renderer").Call("on", "afterRender", func() {
		v.Resize()
	})

	var last *struct{}
	v.editor.OnChange(func(ev *js.Object) {
		last = &struct{}{}
		before := last
		go func() {
			<-time.After(time.Millisecond * 250)
			if before == last {
				value := v.editor.GetValue()
				if value == v.app.Source.Current() {
					// don't fire event if text hasn't changed
					return
				}
				v.app.Dispatch(&actions.UserChangedText{
					Text: value,
				})
			}
		}()
	})
}

func (v *Editor) Resize() {
	v.editor.Call("resize")
}

func (v *Editor) Unmount() {
	v.app.Delete(v)
}

func (v *Editor) Render() vecty.ComponentOrHTML {

	editorDisplay := "none"
	if len(v.app.Source.Packages()) > 0 && len(v.app.Source.Files(v.app.Editor.CurrentPackage())) > 0 {
		editorDisplay = ""
	}

	return elem.Div(
		vecty.Markup(
			prop.ID("editor"),
			vecty.Class("editor"),
			vecty.Style("display", editorDisplay),
		),
	)
}
