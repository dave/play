package views

import (
	"github.com/dave/play/actions"
	"github.com/dave/play/stores"
	"github.com/dave/splitter"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
)

type Page struct {
	vecty.Core
	app *stores.App

	Sizes []float64 `vecty:"prop"`

	split  *splitter.Split
	editor *Editor
}

func NewPage(app *stores.App) *Page {
	v := &Page{
		app: app,
	}
	return v
}

func (v *Page) Mount() {
	v.app.Watch(v, func(done chan struct{}) {
		defer close(done)
		v.Sizes = v.app.Editor.Sizes()
		v.split.SetSizesIfChanged(v.Sizes)

		// Only top-level page should fire vecty.Rerender
		vecty.Rerender(v)
	})

	v.split = splitter.New("split")
	v.split.Init(
		js.S{"#left", "#right"},
		js.M{
			"sizes": v.Sizes,
			"onDragEnd": func() {
				v.editor.Resize()
				v.app.Dispatch(&actions.UserChangedSplit{
					Sizes: v.split.GetSizes(),
				})
			},
		},
	)
}

func (v *Page) Unmount() {
	v.app.Delete(v)
}

const Styles = `
	html, body {
		height: 100%;
	}
	#left {
		display: flex;
		flex-flow: column;
		height: 100%;
	}
	.menu {
		min-height: 56px;
	}
	.editor {
		flex: 1;
		width: 100%;
	}
	.split {
		height: 100%;
		width: 100%;
	}
	.gutter {
		height: 100%;
		background-color: #eee;
		background-repeat: no-repeat;
		background-position: 50%;
	}
	.gutter.gutter-horizontal {
		cursor: col-resize;
		background-image:  url('data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAUAAAAeCAYAAADkftS9AAAAIklEQVQoU2M4c+bMfxAGAgYYmwGrIIiDjrELjpo5aiZeMwF+yNnOs5KSvgAAAABJRU5ErkJggg==')
	}
	.split {
		-webkit-box-sizing: border-box;
		-moz-box-sizing: border-box;
		box-sizing: border-box;
	}
	.split, .gutter.gutter-horizontal {
		float: left;
	}
	.preview {
		border: 0;
		height: 100%;
		width: 100%;
	}
`

func (v *Page) Render() vecty.ComponentOrHTML {
	return elem.Body(
		elem.Div(
			vecty.Markup(
				vecty.Class("container-fluid", "p-0", "split", "split-horizontal"),
			),
			v.renderLeft(),
			v.renderRight(),
			NewAddFileModal(v.app),
			NewDeleteFileModal(v.app),
		),
	)
}

func (v *Page) renderLeft() *vecty.HTML {

	v.editor = NewEditor(v.app)

	return elem.Div(
		vecty.Markup(
			prop.ID("left"),
			vecty.Class("split"),
		),
		NewMenu(v.app),
		v.editor,
	)
}

func (v *Page) renderRight() *vecty.HTML {
	return elem.Div(
		vecty.Markup(
			prop.ID("right"),
			vecty.Class("split"),
		),
		elem.Div(
			vecty.Markup(
				prop.ID("iframe-holder"),
			),
		),
	)
}
