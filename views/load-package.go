package views

import (
	"github.com/dave/play/actions"
	"github.com/dave/play/models"
	"github.com/dave/play/stores"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
)

type LoadPackageModal struct {
	*Modal
	imps *vecty.HTML
}

func NewLoadPackageModal(app *stores.App) *LoadPackageModal {
	v := &LoadPackageModal{}
	v.Modal = &Modal{
		app:    app,
		id:     models.LoadPackageModal,
		title:  "Load package...",
		action: v.action,
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

	return v.Body(
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
	v.app.Dispatch(&actions.ModalClose{Modal: models.LoadPackageModal})
	v.app.Dispatch(&actions.GetStart{Path: value})
}
