package views

import (
	"strings"

	"github.com/dave/play/actions"
	"github.com/dave/play/models"
	"github.com/dave/play/stores"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
)

type BuildTagsModal struct {
	*Modal
	input *vecty.HTML
}

func NewBuildTagsModal(app *stores.App) *BuildTagsModal {
	v := &BuildTagsModal{}
	v.Modal = &Modal{
		app:    app,
		id:     models.BuildTagsModal,
		title:  "Build tags",
		action: v.action,
	}
	return v
}

func (v *BuildTagsModal) Render() vecty.ComponentOrHTML {
	v.input = elem.Input(vecty.Markup(
		vecty.Class("form-control"),
		prop.Type(prop.TypeText),
		prop.ID("build-tags-input"),
		prop.Value(strings.Join(v.app.Compile.Tags(), " ")),
	))

	return v.Body(
		elem.Form(
			elem.Div(
				vecty.Markup(
					vecty.Class("form-group"),
				),
				elem.Label(
					vecty.Markup(
						vecty.Property("for", "build-tags-input"),
						vecty.Class("col-form-label"),
					),
					vecty.Text("Build tags"),
				),
				v.input,
			),
		),
	).Build()
}

func (v *BuildTagsModal) action(*vecty.Event) {
	tags := strings.Fields(v.input.Node().Get("value").String())
	v.app.Dispatch(&actions.ModalClose{Modal: models.BuildTagsModal})
	if !compare(v.app.Compile.Tags(), tags) {
		v.app.Dispatch(&actions.BuildTags{Tags: tags})
	}
}

func compare(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
