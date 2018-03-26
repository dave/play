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

	Text string `vecty:"prop"`

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
		v.Text = v.app.Editor.Text()
		if v.Text != v.editor.GetValue() {
			// only update the editor if the text is changed
			v.editor.SetValue(v.Text)
			v.editor.ClearSelection()
			v.editor.MoveCursorTo(0, 0)
		}
	})

	v.editor = ace.Edit("editor")
	v.editor.SetOptions(map[string]interface{}{
		"mode": "ace/mode/golang",
	})
	if v.Text != "" {
		v.editor.SetValue(v.Text)
		v.editor.ClearSelection()
		v.editor.MoveCursorTo(0, 0)
	}

	dom.GetWindow().AddEventListener("resize", false, func(event dom.Event) {
		v.Resize()
	})

	var changes int
	v.editor.OnChange(func(ev *js.Object) {
		changes++
		before := changes
		go func() {
			<-time.After(time.Millisecond * 250)
			if before == changes {
				value := v.editor.GetValue()
				if value == v.app.Editor.Files()[v.app.Editor.Current()] {
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
	return elem.Div(
		vecty.Markup(
			prop.ID("editor"),
			vecty.Class("editor"),
		),
	)
}
