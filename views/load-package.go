package views

import (
	"sort"

	"github.com/dave/play/actions"
	"github.com/dave/play/models"
	"github.com/dave/play/stores"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
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
		title:  "Load package",
		action: v.action,
		hidden: func() {
			app.Dispatch(&actions.ShowAllDepsChange{State: false})
		},
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
	var paths []string
	if !v.app.Page.ShowAllDeps() {
		paths = v.app.Scanner.AllImportsOrdered()
	} else {
		imps := v.app.Scanner.AllImports()
		for p := range imps {
			paths = append(paths, p)
		}
		for p := range v.app.Archive.Cache() {
			if !imps[p] {
				paths = append(paths, p)
			}
		}
		sort.Strings(paths)
	}
	for _, path := range paths {
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

	infoDisplay := "none"
	if v.app.Page.ShowAllDeps() && !v.app.Archive.AllFresh() {
		infoDisplay = ""
	}

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
			elem.Div(
				vecty.Markup(
					vecty.Class("form-check"),
				),
				elem.Input(
					vecty.Markup(
						prop.Type(prop.TypeCheckbox),
						vecty.Class("form-check-input"),
						prop.ID("load-package-show-all-deps-checkbox"),
						prop.Checked(v.app.Page.ShowAllDeps()),
						event.Click(func(e *vecty.Event) {
							v.app.Dispatch(&actions.ShowAllDepsChange{
								State: e.Target.Get("checked").Bool(),
							})
						}).PreventDefault(),
					),
				),
				elem.Label(
					vecty.Markup(
						vecty.Class("form-check-label"),
						prop.For("load-package-show-all-deps-checkbox"),
					),
					vecty.Text("Show all dependencies"),
				),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("form-text"),
					vecty.Style("display", infoDisplay),
				),
				vecty.Tag(
					"svg",
					vecty.Markup(
						vecty.Namespace("http://www.w3.org/2000/svg"),
						vecty.Attribute("width", "14"),
						vecty.Attribute("height", "16"),
						vecty.Attribute("viewBox", "0 0 14 16"),
						vecty.Class("octicon"),
						vecty.Style("margin-right", "4px"),
					),
					vecty.Tag(
						"path",
						vecty.Markup(
							vecty.Namespace("http://www.w3.org/2000/svg"),
							vecty.Attribute("fill-rule", "evenodd"),
							vecty.Attribute("d", "M6.3 5.71a.942.942 0 0 1-.28-.7c0-.28.09-.52.28-.7.19-.18.42-.28.7-.28.28 0 .52.09.7.28.18.19.28.42.28.7 0 .28-.09.52-.28.7a1 1 0 0 1-.7.3c-.28 0-.52-.11-.7-.3zM8 8.01c-.02-.25-.11-.48-.31-.69-.2-.19-.42-.3-.69-.31H6c-.27.02-.48.13-.69.31-.2.2-.3.44-.31.69h1v3c.02.27.11.5.31.69.2.2.42.31.69.31h1c.27 0 .48-.11.69-.31.2-.19.3-.42.31-.69H8V8v.01zM7 2.32C3.86 2.32 1.3 4.86 1.3 8c0 3.14 2.56 5.7 5.7 5.7s5.7-2.55 5.7-5.7c0-3.15-2.56-5.69-5.7-5.69v.01zM7 1c3.86 0 7 3.14 7 7s-3.14 7-7 7-7-3.12-7-7 3.14-7 7-7z"),
						),
					),
				),
				vecty.Text("Dependencies are not fully loaded. "),
				elem.Anchor(
					vecty.Markup(
						prop.Href(""),
						event.Click(func(e *vecty.Event) {
							v.app.Dispatch(&actions.ModalClose{Modal: models.LoadPackageModal})
							v.app.Dispatch(&actions.UpdateStart{})
						}).PreventDefault(),
					),
					vecty.Text("Update"),
				),
				vecty.Text(" to load."),
			),
		),
	).Build()
}

func (v *LoadPackageModal) action(*vecty.Event) {
	n := v.imps.Node()
	i := n.Get("selectedIndex").Int()
	value := n.Get("options").Index(i).Get("value").String()
	v.app.Dispatch(&actions.ModalClose{Modal: models.LoadPackageModal})
	v.app.Dispatch(&actions.GetStart{Path: value, Save: true})
}
