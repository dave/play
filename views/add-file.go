package views

import (
	"fmt"

	"strings"

	"github.com/dave/play/actions"
	"github.com/dave/play/models"
	"github.com/dave/play/stores"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
)

type AddFileModal struct {
	*Modal
	input *vecty.HTML
}

func NewAddFileModal(app *stores.App) *AddFileModal {
	v := &AddFileModal{}
	v.Modal = &Modal{
		app:    app,
		id:     models.AddFileModal,
		title:  "Add file",
		action: v.save,
		shown: func() {
			js.Global.Call("$", "#add-file-input").Call("focus")
			js.Global.Call("$", "#add-file-input").Call("val", "")
		},
	}
	return v
}

func (v *AddFileModal) Render() vecty.ComponentOrHTML {
	v.input = elem.Input(
		vecty.Markup(
			prop.Type(prop.TypeText),
			vecty.Class("form-control"),
			prop.ID("add-file-input"),
			event.KeyPress(func(ev *vecty.Event) {
				if ev.Get("keyCode").Int() == 13 {
					ev.Call("preventDefault")
					v.save(ev)
				}
			}),
		),
	)
	return v.Body(
		elem.Form(
			elem.Div(
				vecty.Markup(vecty.Class("form-group")),
				elem.Label(
					vecty.Markup(
						vecty.Property("for", "add-file-input"),
						vecty.Class("col-form-label"),
					),
					vecty.Text("Filename"),
				),
				v.input,
			),
		),
	).Build()
}

func (v *AddFileModal) save(*vecty.Event) {
	value := v.input.Node().Get("value").String()
	if strings.Contains(value, "/") {
		v.app.Fail(fmt.Errorf("filename %s must not contain a slash", value))
		return
	}
	if !strings.HasSuffix(value, ".go") && !strings.Contains(value, ".") {
		value = value + ".go"
	}
	if v.app.Source.HasFile(v.app.Editor.CurrentPackage(), value) {
		v.app.Fail(fmt.Errorf("%s already exists", value))
		return
	}
	v.app.Dispatch(&actions.ModalClose{Modal: models.AddFileModal})
	v.app.Dispatch(&actions.AddFile{Name: value})
}
