package views

import (
	"github.com/dave/play/models"
	"github.com/dave/play/stores"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
)

type ClashWarningModal struct {
	*Modal
}

func NewClashWarningModal(app *stores.App) *ClashWarningModal {
	v := &ClashWarningModal{
		&Modal{
			app:    app,
			id:     models.ClashWarningModal,
			title:  "Package clash...",
			action: nil,
		},
	}
	return v
}

func (v *ClashWarningModal) Render() vecty.ComponentOrHTML {
	return v.Body(
		elem.Form(
			elem.Div(
				vecty.Markup(vecty.Class("form-group")),
				elem.Label(
					vecty.Markup(
						vecty.Property("for", string(v.id)),
						vecty.Class("col-form-label"),
					),
					vecty.Text("Link"),
				),
				elem.Input(
					vecty.Markup(
						prop.Type(prop.TypeText),
						vecty.Class("form-control"),
						prop.ID("clash-warning-input-link"),
						event.Focus(func(ev *vecty.Event) {
							ev.Target.Call("select")
						}).PreventDefault(),
						prop.Value(v.app.Deploy.Index()),
					),
				),
				elem.Small(
					vecty.Markup(
						vecty.Class("form-text", "text-muted"),
					),
					elem.Anchor(
						vecty.Markup(
							prop.Href(v.app.Deploy.Index()),
							vecty.Property("target", "_blank"),
						),
						vecty.Text("Click here"),
					),
					vecty.Text(". Use the link for testing and toy projects. Remember you're sharing the jsgo.io domain with everyone else, so the browser environment should be considered toxic."),
				),
				elem.Label(
					vecty.Markup(
						vecty.Property("for", string(v.id)),
						vecty.Class("col-form-label"),
					),
					vecty.Text("Loader JS"),
				),
				elem.Input(
					vecty.Markup(
						prop.Type(prop.TypeText),
						vecty.Class("form-control"),
						prop.ID("clash-warning-input-loader"),
						event.Focus(func(ev *vecty.Event) {
							ev.Target.Call("select")
						}).PreventDefault(),
						prop.Value(v.app.Deploy.LoaderJs()),
					),
				),
				elem.Small(
					vecty.Markup(
						vecty.Class("form-text", "text-muted"),
					),
					vecty.Text("For production, use the Loader JS in a script tag on your own site."),
				),
			),
		),
	).Build()
}
