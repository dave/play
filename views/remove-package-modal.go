package views

import (
	"github.com/dave/play/actions"
	"github.com/dave/play/stores"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
)

type RemovePackageModal struct {
	vecty.Core
	app *stores.App
	sel *vecty.HTML
}

func NewRemovePackageModal(app *stores.App) *RemovePackageModal {
	v := &RemovePackageModal{
		app: app,
	}
	return v
}

func (v *RemovePackageModal) Render() vecty.ComponentOrHTML {
	items := []vecty.MarkupOrChild{
		vecty.Markup(
			vecty.Class("form-control"),
			prop.ID("remove-package-select"),
		),
	}
	for _, path := range v.app.Source.Packages() {
		items = append(items,
			elem.Option(
				vecty.Markup(
					prop.Value(path),
					vecty.Property("selected", v.app.Editor.CurrentPackage() == path),
				),
				vecty.Text(path),
			),
		)
	}
	v.sel = elem.Select(items...)

	return Modal(
		"Remove package...",
		"remove-package-modal",
		v.action,
	).Body(
		elem.Form(
			elem.Div(
				vecty.Markup(
					vecty.Class("form-group"),
				),
				elem.Label(
					vecty.Markup(
						vecty.Property("for", "remove-package-select"),
						vecty.Class("col-form-label"),
					),
					vecty.Text("Package path"),
				),
				v.sel,
			),
		),
	).Build()
}

func (v *RemovePackageModal) action(*vecty.Event) {
	n := v.sel.Node()
	i := n.Get("selectedIndex").Int()
	value := n.Get("options").Index(i).Get("value").String()
	v.app.Dispatch(&actions.RemovePackage{
		Path: value,
	})
}
