package views

import (
	"github.com/dave/play/actions"
	"github.com/dave/play/models"
	"github.com/dave/play/stores"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
)

type DeleteFileModal struct {
	*Modal
	sel *vecty.HTML
}

func NewDeleteFileModal(app *stores.App) *DeleteFileModal {
	v := &DeleteFileModal{}
	v.Modal = &Modal{
		app:    app,
		id:     models.DeleteFileModal,
		title:  "Delete file",
		action: v.action,
	}
	return v
}

func (v *DeleteFileModal) Render() vecty.ComponentOrHTML {
	items := []vecty.MarkupOrChild{
		vecty.Markup(
			vecty.Class("form-control"),
			prop.ID("delete-file-select"),
		),
	}
	for _, name := range v.app.Source.Filenames(v.app.Editor.CurrentPackage()) {
		items = append(items,
			elem.Option(
				vecty.Markup(
					prop.Value(name),
					vecty.Property("selected", v.app.Editor.CurrentFile() == name),
				),
				vecty.Text(name),
			),
		)
	}
	v.sel = elem.Select(items...)

	return v.Body(
		elem.Form(
			elem.Div(
				vecty.Markup(
					vecty.Class("form-group"),
				),
				elem.Label(
					vecty.Markup(
						vecty.Property("for", "delete-file-select"),
						vecty.Class("col-form-label"),
					),
					vecty.Text("File"),
				),
				v.sel,
			),
		),
	).Build()
}

func (v *DeleteFileModal) action(*vecty.Event) {
	n := v.sel.Node()
	i := n.Get("selectedIndex").Int()
	value := n.Get("options").Index(i).Get("value").String()
	v.app.Dispatch(&actions.ModalClose{Modal: models.DeleteFileModal})
	v.app.Dispatch(&actions.DeleteFile{Name: value})
}
