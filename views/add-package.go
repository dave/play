package views

import (
	"fmt"

	"github.com/dave/play/actions"
	"github.com/dave/play/models"
	"github.com/dave/play/stores"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
)

type AddPackageModal struct {
	*Modal
	input *vecty.HTML
}

func NewAddPackageModal(app *stores.App) *AddPackageModal {
	v := &AddPackageModal{}
	v.Modal = &Modal{
		app:    app,
		id:     models.AddPackageModal,
		title:  "Add package...",
		action: v.save,
		shown: func() {
			js.Global.Call("$", "#add-package-input").Call("focus")
			js.Global.Call("$", "#add-package-input").Call("val", "")
		},
	}
	return v
}

func (v *AddPackageModal) Render() vecty.ComponentOrHTML {
	v.input = elem.Input(
		vecty.Markup(
			prop.Type(prop.TypeText),
			vecty.Class("form-control"),
			prop.ID("add-package-input"),
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
						vecty.Property("for", "add-package-input"),
						vecty.Class("col-form-label"),
					),
					vecty.Text("Package path"),
				),
				v.input,
			),
		),
	).Build()
}

func (v *AddPackageModal) save(*vecty.Event) {
	value := v.input.Node().Get("value").String()
	if v.app.Source.HasPackage(value) {
		v.app.Fail(fmt.Errorf("%s already exists", value))
		return
	}
	v.app.Dispatch(&actions.ModalClose{Modal: models.AddPackageModal})
	v.app.Dispatch(&actions.AddPackage{Path: value})
}
