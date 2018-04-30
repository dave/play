package views

import (
	"github.com/dave/dropper"
	"github.com/dave/play/actions"
	"github.com/dave/play/models"
	"github.com/dave/play/stores"
	"github.com/dave/splitter"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
	"honnef.co/go/js/dom"
)

type Page struct {
	vecty.Core
	app *stores.App

	split1, split2 *splitter.Split
	editor         *Editor
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

		sizes := v.app.Editor.Sizes()
		if v.split1.Changed(sizes) {
			v.split1.SetSizes(sizes)
		}

		if v.app.Page.Console() && !v.split2.Initialised() {
			v.split2.Init(
				js.S{"#iframe-holder", "#console-holder"},
				js.M{
					"direction": "vertical",
					"sizes":     []float64{50.0, 50.0},
				},
			)
		} else if !v.app.Page.Console() && v.split2.Initialised() {
			v.split2.Destroy()
		}
	})

	v.split1 = splitter.New("split")
	v.split1.Init(
		js.S{"#left", "#right"},
		js.M{
			"sizes": v.app.Editor.Sizes(),
			"onDragEnd": func() {
				v.editor.Resize()
				v.app.Dispatch(&actions.UserChangedSplit{
					Sizes: v.split1.GetSizes(),
				})
			},
		},
	)

	v.split2 = splitter.New("split")

	enter, leave, drop := dropper.Initialise(dom.GetWindow().Document().GetElementByID("left"))
	go func() {
		for {
			select {
			case <-enter:
				v.app.Dispatch(&actions.DragEnter{})
			case <-leave:
				v.app.Dispatch(&actions.DragLeave{})
			case files := <-drop:
				v.app.Dispatch(&actions.DragDrop{
					Files: files,
				})
			}
		}
	}()

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
	.editor, .empty-panel {
		flex: 1;
		width: 100%;
	}
	.empty-panel {
		display: flex;
		align-items: center;
		justify-content: center;
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
	.gutter.gutter-vertical {
		cursor: row-resize;
		background-image:  url('data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAB4AAAAFAQMAAABo7865AAAABlBMVEVHcEzMzMzyAv2sAAAAAXRSTlMAQObYZgAAABBJREFUeF5jOAMEEAIEEFwAn3kMwcB6I2AAAAAASUVORK5CYII=')
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
	#console-holder {
		overflow: auto;
	}
	#console {
		padding:5px;
	}
	.octicon {
		display: inline-block;
		vertical-align: text-top;
		fill: currentColor;
	}
	#help-modal table { 
		clear: both;
	}
	#help-modal img { 
		margin-left: 20px;
		margin-bottom: 30px;
	}
`

func (v *Page) Render() vecty.ComponentOrHTML {
	githubBannerDisplay := ""
	if v.app.Compile.Compiled() {
		githubBannerDisplay = "none"
	}
	return elem.Body(
		elem.Div(
			vecty.Markup(
				vecty.Class("container-fluid", "p-0", "split", "split-horizontal"),
			),
			v.renderLeft(),
			v.renderRight(),
		),
		NewAddFileModal(v.app),
		NewDeleteFileModal(v.app),
		NewAddPackageModal(v.app),
		NewRemovePackageModal(v.app),
		NewDeployDoneModal(v.app),
		NewLoadPackageModal(v.app),
		NewClashWarningModal(v.app),
		NewBuildTagsModal(v.app),
		NewHelpModal(v.app),
		elem.Anchor(
			vecty.Markup(
				prop.Href("https://github.com/dave/play"),
				vecty.Style("display", githubBannerDisplay),
				vecty.Property("target", "_blank"),
			),
			elem.Image(
				vecty.Markup(
					vecty.Style("position", "absolute"),
					vecty.Style("top", "0"),
					vecty.Style("right", "0"),
					vecty.Style("border", "0"),
					prop.Src("https://s3.amazonaws.com/github/ribbons/forkme_right_gray_6d6d6d.png"),
					vecty.Property("alt", "Fork me on GitHub"),
				),
			),
		),
	)
}

func (v *Page) renderLeft() *vecty.HTML {

	v.editor = NewEditor(v.app)

	emptyDisplay := "none"
	addFileDisplay := "none"
	addPackageDisplay := "none"
	loadingDisplay := "none"
	if !v.app.Editor.Loaded() {
		emptyDisplay = ""
		loadingDisplay = ""
	} else if len(v.app.Source.Packages()) == 0 {
		emptyDisplay = ""
		addPackageDisplay = ""
	} else if len(v.app.Source.Files(v.app.Editor.CurrentPackage())) == 0 {
		emptyDisplay = ""
		addFileDisplay = ""
	}

	return elem.Div(
		vecty.Markup(
			prop.ID("left"),
			vecty.Class("split"),
		),
		NewMenu(v.app),
		v.editor,
		elem.Div(
			vecty.Markup(
				vecty.Class("empty-panel"),
				vecty.Style("display", emptyDisplay),
			),
			elem.Span(
				vecty.Markup(
					vecty.Style("display", loadingDisplay),
				),
				vecty.Text("Loading..."),
			),
			elem.Button(
				vecty.Markup(
					vecty.Property("type", "button"),
					vecty.Class("btn", "btn-primary"),
					event.Click(func(e *vecty.Event) {
						v.app.Dispatch(&actions.ModalOpen{Modal: models.AddFileModal})
					}).PreventDefault(),
					vecty.Style("display", addFileDisplay),
				),
				vecty.Text("Add file"),
			),
			elem.Button(
				vecty.Markup(
					vecty.Property("type", "button"),
					vecty.Class("btn", "btn-primary"),
					event.Click(func(e *vecty.Event) {
						v.app.Dispatch(&actions.ModalOpen{Modal: models.AddPackageModal})
					}).PreventDefault(),
					vecty.Style("display", addPackageDisplay),
				),
				vecty.Text("Add package"),
			),
		),
	)
}

func (v *Page) renderRight() *vecty.HTML {
	consoleDisplay := ""
	if !v.app.Page.Console() {
		consoleDisplay = "none"
	}
	return elem.Div(
		vecty.Markup(
			prop.ID("right"),
			vecty.Class("split", "split-vertical"),
		),
		elem.Div(
			vecty.Markup(
				prop.ID("iframe-holder"),
			),
		),
		elem.Div(
			vecty.Markup(
				prop.ID("console-holder"),
				vecty.Style("display", consoleDisplay),
			),
			elem.Preformatted(
				vecty.Markup(
					prop.ID("console"),
				),
			),
		),
	)
}
