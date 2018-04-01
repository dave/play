package views

import (
	"fmt"

	"github.com/dave/play/actions"
	"github.com/dave/play/stores"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
)

type AddPackageModal struct {
	vecty.Core
	app   *stores.App
	input *vecty.HTML
}

func NewAddPackageModal(app *stores.App) *AddPackageModal {
	v := &AddPackageModal{
		app: app,
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
	return Modal(
		"Add package...",
		"add-package-modal",
		v.save,
	).Body(
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
	v.app.Dispatch(&actions.AddPackage{
		Path: value,
	})
}
