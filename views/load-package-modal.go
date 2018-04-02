package views

import (
	"github.com/dave/play/actions"
	"github.com/dave/play/stores"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
)

type LoadPackageModal struct {
	vecty.Core
	app  *stores.App
	imps *vecty.HTML
}

func NewLoadPackageModal(app *stores.App) *LoadPackageModal {
	v := &LoadPackageModal{
		app: app,
	}
	return v
}

func (v *LoadPackageModal) Render() vecty.ComponentOrHTML {
	items := []vecty.MarkupOrChild{
		vecty.Markup(
			vecty.Class("form-control"),
			prop.ID("load-package-imps-select"),
		),
	}
	for _, path := range v.app.Scanner.AllImportsOrdered() {
		items = append(items,
			elem.Option(
				vecty.Markup(
					prop.Value(path),
				),
				vecty.Text(path),
			),
		)
	}
	v.imps = elem.Select(items...)

	return Modal(
		"Load package...",
		"load-package-modal",
		v.action,
	).Body(
		elem.Form(
			elem.Div(
				vecty.Markup(
					vecty.Class("form-group"),
				),
				elem.Label(
					vecty.Markup(
						vecty.Property("for", "load-package-imps-select"),
						vecty.Class("col-form-label"),
					),
					vecty.Text("Imports"),
				),
				v.imps,
			),
		),
	).Build()
}

func (v *LoadPackageModal) action(*vecty.Event) {
	n := v.imps.Node()
	i := n.Get("selectedIndex").Int()
	value := n.Get("options").Index(i).Get("value").String()
	v.app.Dispatch(&actions.GetStart{
		Path: value,
	})
}
