package views

import (
	"fmt"

	"github.com/dave/play/actions"
	"github.com/dave/play/models"
	"github.com/dave/play/stores"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
)

type Menu struct {
	vecty.Core
	app *stores.App

	compileButton *vecty.HTML
	optionsButton *vecty.HTML
}

func NewMenu(app *stores.App) *Menu {
	v := &Menu{
		app: app,
	}
	return v
}

func (v *Menu) Render() vecty.ComponentOrHTML {

	clashWarningDisplay := "none"
	if len(v.app.Scanner.Clashes()) > 0 {
		clashWarningDisplay = ""
	}

	buildTagsText := "Build tags..."
	if len(v.app.Compile.Tags()) > 0 {
		buildTagsText = fmt.Sprintf("Build tags (%d)...", len(v.app.Compile.Tags()))
	}

	return elem.Navigation(
		vecty.Markup(
			vecty.Class("menu", "navbar", "navbar-expand", "navbar-light", "bg-light"),
		),
		elem.UnorderedList(
			vecty.Markup(
				vecty.Class("navbar-nav", "mr-auto"),
			),
			v.renderPackageDropdown(),
			v.renderFileDropdown(),

			elem.ListItem(
				vecty.Markup(
					vecty.Class("nav-item"),
					vecty.Style("display", clashWarningDisplay),
				),
				elem.Anchor(
					vecty.Markup(
						prop.Href(""),
						vecty.Class("nav-link", "octicon"),
						event.Click(func(e *vecty.Event) {
							v.app.Dispatch(&actions.ModalOpen{Modal: models.ClashWarningModal})
						}).PreventDefault(),
					),
					vecty.Tag(
						"svg",
						vecty.Markup(
							vecty.Namespace("http://www.w3.org/2000/svg"),
							vecty.Attribute("width", "14"),
							vecty.Attribute("height", "16"),
							vecty.Attribute("viewBox", "0 0 14 16"),
						),
						vecty.Tag(
							"path",
							vecty.Markup(
								vecty.Namespace("http://www.w3.org/2000/svg"),
								vecty.Attribute("fill-rule", "evenodd"),
								vecty.Attribute("d", "M7 2.3c3.14 0 5.7 2.56 5.7 5.7s-2.56 5.7-5.7 5.7A5.71 5.71 0 0 1 1.3 8c0-3.14 2.56-5.7 5.7-5.7zM7 1C3.14 1 0 4.14 0 8s3.14 7 7 7 7-3.14 7-7-3.14-7-7-7zm1 3H6v5h2V4zm0 6H6v2h2v-2z"),
							),
						),
					),
				),
			),
		),
		elem.UnorderedList(
			vecty.Markup(
				vecty.Class("navbar-nav", "ml-auto"),
			),
			elem.ListItem(
				vecty.Markup(
					vecty.Class("nav-item"),
				),
				elem.Span(
					vecty.Markup(
						vecty.Class("navbar-text"),
						prop.ID("message"),
					),
					vecty.Text(""),
				),
			),
			elem.ListItem(
				vecty.Markup(
					vecty.Class("nav-item"),
				),
				elem.Anchor(
					vecty.Markup(
						prop.Href(""),
						//vecty.Style("margin-right", "5px"),
						vecty.Class("nav-link", "octicon"),
						event.Click(func(e *vecty.Event) {
							v.app.Dispatch(&actions.ModalOpen{Modal: models.HelpModal})
						}).PreventDefault(),
					),
					vecty.Tag(
						"svg",
						vecty.Markup(
							vecty.Namespace("http://www.w3.org/2000/svg"),
							vecty.Attribute("width", "14"),
							vecty.Attribute("height", "16"),
							vecty.Attribute("viewBox", "0 0 14 16"),
						),
						vecty.Tag(
							"path",
							vecty.Markup(
								vecty.Namespace("http://www.w3.org/2000/svg"),
								vecty.Attribute("fill-rule", "evenodd"),
								vecty.Attribute("d", "M6 10h2v2H6v-2zm4-3.5C10 8.64 8 9 8 9H6c0-.55.45-1 1-1h.5c.28 0 .5-.22.5-.5v-1c0-.28-.22-.5-.5-.5h-1c-.28 0-.5.22-.5.5V7H4c0-1.5 1.5-3 3-3s3 1 3 2.5zM7 2.3c3.14 0 5.7 2.56 5.7 5.7s-2.56 5.7-5.7 5.7A5.71 5.71 0 0 1 1.3 8c0-3.14 2.56-5.7 5.7-5.7zM7 1C3.14 1 0 4.14 0 8s3.14 7 7 7 7-3.14 7-7-3.14-7-7-7z"),
							),
						),
					),
				),
			),
			elem.ListItem(
				vecty.Markup(
					vecty.Class("nav-item", "btn-group"),
				),
				elem.Button(
					vecty.Markup(
						vecty.Property("type", "button"),
						vecty.Class("btn", "btn-primary"),
						event.Click(func(e *vecty.Event) {
							if v.app.Connection.Open() || v.app.Compile.Compiling() {
								return
							} else {
								v.app.Dispatch(&actions.FormatCode{
									Then: &actions.CompileStart{},
								})
							}
						}).PreventDefault(),
					),
					vecty.Text("Run"),
				),
				elem.Button(
					vecty.Markup(
						vecty.Property("type", "button"),
						vecty.Data("toggle", "dropdown"),
						vecty.Property("aria-haspopup", "true"),
						vecty.Property("aria-expanded", "false"),
						vecty.Class("btn", "btn-primary", "dropdown-toggle", "dropdown-toggle-split"),
					),
					elem.Span(vecty.Markup(vecty.Class("sr-only")), vecty.Text("Options")),
				),
				elem.Div(
					vecty.Markup(
						vecty.Class("dropdown-menu", "dropdown-menu-right"),
					),
					elem.Anchor(
						vecty.Markup(
							vecty.Class("dropdown-item"),
							prop.Href(""),
							event.Click(func(e *vecty.Event) {
								v.app.Dispatch(&actions.FormatCode{})
							}).PreventDefault(),
						),
						vecty.Text("Format code"),
					),
					elem.Div(
						vecty.Markup(
							vecty.Class("dropdown-divider"),
						),
					),
					elem.Anchor(
						vecty.Markup(
							vecty.Class("dropdown-item"),
							prop.Href(""),
							event.Click(func(e *vecty.Event) {
								v.app.Dispatch(&actions.FormatCode{
									Then: &actions.RequestStart{Type: models.UpdateRequest},
								})
							}).PreventDefault(),
						),
						vecty.Text("Update"),
					),
					elem.Anchor(
						vecty.Markup(
							vecty.Class("dropdown-item"),
							prop.Href(""),
							event.Click(func(e *vecty.Event) {
								v.app.Dispatch(&actions.FormatCode{
									Then: &actions.ShareStart{},
								})
							}).PreventDefault(),
						),
						vecty.Text("Share"),
					),
					elem.Anchor(
						vecty.Markup(
							vecty.Class("dropdown-item"),
							prop.Href(""),
							event.Click(func(e *vecty.Event) {
								v.app.Dispatch(&actions.FormatCode{
									Then: &actions.DeployStart{},
								})
							}).PreventDefault(),
						),
						vecty.Text("Deploy"),
					),
					elem.Div(
						vecty.Markup(
							vecty.Class("dropdown-divider"),
						),
					),
					elem.Anchor(
						vecty.Markup(
							vecty.Class("dropdown-item"),
							prop.Href("#"),
							event.Click(func(e *vecty.Event) {}).StopPropagation(),
						),
						elem.Input(
							vecty.Markup(
								prop.Type(prop.TypeCheckbox),
								vecty.Class("form-check-input", "dropdown-item"),
								prop.ID("dropdownCheckConsole"),
								prop.Checked(v.app.Page.Console()),
								event.Change(func(e *vecty.Event) {
									v.app.Dispatch(&actions.ConsoleToggleClick{})
								}),
								vecty.Style("cursor", "pointer"),
							),
						),
						elem.Label(
							vecty.Markup(
								vecty.Class("form-check-label"),
								prop.For("dropdownCheckConsole"),
								vecty.Style("cursor", "pointer"),
							),
							vecty.Text("Show console"),
						),
					),
					elem.Anchor(
						vecty.Markup(
							vecty.Class("dropdown-item"),
							prop.Href("#"),
							event.Click(func(e *vecty.Event) {}).StopPropagation(),
						),
						elem.Input(
							vecty.Markup(
								prop.Type(prop.TypeCheckbox),
								vecty.Class("form-check-input", "dropdown-item"),
								prop.ID("dropdownCheckMinify"),
								prop.Checked(v.app.Page.Minify()),
								event.Change(func(e *vecty.Event) {
									v.app.Dispatch(&actions.MinifyToggleClick{})
								}),
								vecty.Style("cursor", "pointer"),
							),
						),
						elem.Label(
							vecty.Markup(
								vecty.Class("form-check-label"),
								prop.For("dropdownCheckMinify"),
								vecty.Style("cursor", "pointer"),
							),
							vecty.Text("Minify JS"),
						),
					),
					elem.Anchor(
						vecty.Markup(
							vecty.Class("dropdown-item"),
							prop.Href(""),
							event.Click(func(e *vecty.Event) {
								v.app.Dispatch(&actions.ModalOpen{Modal: models.BuildTagsModal})
							}).PreventDefault(),
						),
						vecty.Text(buildTagsText),
					),
					elem.Div(
						vecty.Markup(
							vecty.Class("dropdown-divider"),
						),
					),
					elem.Anchor(
						vecty.Markup(
							vecty.Class("dropdown-item"),
							prop.Href(""),
							event.Click(func(e *vecty.Event) {
								v.app.Dispatch(&actions.DownloadClick{})
							}).PreventDefault(),
						),
						vecty.Text("Download"),
					),
					elem.Div(
						vecty.Markup(
							vecty.Class("dropdown-divider"),
						),
					),
					elem.Anchor(
						vecty.Markup(
							vecty.Class("dropdown-item"),
							prop.Href("https://github.com/dave/play"),
							vecty.Property("target", "_blank"),
						),
						vecty.Text("More info"),
					),
				),
			),
		),
	)
}

func (v *Menu) renderFileDropdown() *vecty.HTML {
	var fileItems []vecty.MarkupOrChild
	fileItems = append(fileItems,
		vecty.Markup(
			vecty.Class("dropdown-menu"),
			vecty.Property("aria-labelledby", "fileDropdown"),
		),
	)
	for _, name := range v.app.Source.Filenames(v.app.Editor.CurrentPackage()) {
		name := name
		fileItems = append(fileItems,
			elem.Anchor(
				vecty.Markup(
					vecty.Class("dropdown-item"),
					vecty.ClassMap{
						"disabled": name == v.app.Editor.CurrentFile(),
					},
					prop.Href(""),
					event.Click(func(e *vecty.Event) {
						v.app.Dispatch(&actions.UserChangedFile{
							Name: name,
						})
					}).PreventDefault(),
				),
				vecty.Text(name),
			),
		)
	}
	fileItems = append(fileItems,
		elem.Div(
			vecty.Markup(
				vecty.Class("dropdown-divider"),
			),
		),
		elem.Anchor(
			vecty.Markup(
				vecty.Class("dropdown-item"),
				prop.Href(""),
				event.Click(func(e *vecty.Event) {
					v.app.Dispatch(&actions.ModalOpen{Modal: models.AddFileModal})
				}).PreventDefault(),
			),
			vecty.Text("Add file"),
		),
		elem.Anchor(
			vecty.Markup(
				vecty.Class("dropdown-item"),
				prop.Href(""),
				event.Click(func(e *vecty.Event) {
					v.app.Dispatch(&actions.ModalOpen{Modal: models.DeleteFileModal})
				}).PreventDefault(),
			),
			vecty.Text("Delete file"),
		),
	)

	classes := vecty.Class("nav-item", "dropdown", "d-none")
	if len(v.app.Source.Files(v.app.Editor.CurrentPackage())) > 0 {
		classes = vecty.Class("nav-item", "dropdown")
	}

	return elem.ListItem(
		vecty.Markup(
			classes,
		),
		elem.Anchor(
			vecty.Markup(
				prop.ID("fileDropdown"),
				prop.Href(""),
				vecty.Class("nav-link", "dropdown-toggle"),
				vecty.Property("role", "button"),
				vecty.Data("toggle", "dropdown"),
				vecty.Property("aria-haspopup", "true"),
				vecty.Property("aria-expanded", "false"),
				event.Click(func(ev *vecty.Event) {}).PreventDefault(),
			),
			vecty.Text(v.app.Editor.CurrentFile()),
		),
		elem.Div(
			fileItems...,
		),
	)
}

func (v *Menu) renderPackageDropdown() *vecty.HTML {
	var packageItems []vecty.MarkupOrChild
	packageItems = append(packageItems,
		vecty.Markup(
			vecty.Class("dropdown-menu"),
			vecty.Property("aria-labelledby", "packageDropdown"),
		),
	)
	for _, path := range v.app.Source.Packages() {
		path := path
		packageItems = append(packageItems,
			elem.Anchor(
				vecty.Markup(
					vecty.Class("dropdown-item"),
					vecty.ClassMap{
						"disabled": path == v.app.Editor.CurrentPackage(),
					},
					prop.Href(""),
					event.Click(func(e *vecty.Event) {
						v.app.Dispatch(&actions.UserChangedPackage{
							Path: path,
						})
					}).PreventDefault(),
				),
				vecty.Text(v.app.Scanner.DisplayPath(path)),
			),
		)
	}
	packageItems = append(packageItems,
		elem.Div(
			vecty.Markup(
				vecty.Class("dropdown-divider"),
			),
		),
		elem.Anchor(
			vecty.Markup(
				vecty.Class("dropdown-item"),
				prop.Href(""),
				event.Click(func(e *vecty.Event) {
					v.app.Dispatch(&actions.ModalOpen{Modal: models.AddPackageModal})
				}).PreventDefault(),
			),
			vecty.Text("Add package"),
		),
		elem.Anchor(
			vecty.Markup(
				vecty.Class("dropdown-item"),
				prop.Href(""),
				event.Click(func(e *vecty.Event) {
					v.app.Dispatch(&actions.ModalOpen{Modal: models.LoadPackageModal})
				}).PreventDefault(),
			),
			vecty.Text("Load package"),
		),
		elem.Anchor(
			vecty.Markup(
				vecty.Class("dropdown-item"),
				prop.Href(""),
				event.Click(func(e *vecty.Event) {
					v.app.Dispatch(&actions.ModalOpen{Modal: models.RemovePackageModal})
				}).PreventDefault(),
			),
			vecty.Text("Remove package"),
		),
	)

	classes := vecty.Class("nav-item", "dropdown", "d-none")
	if len(v.app.Source.Packages()) > 0 {
		classes = vecty.Class("nav-item", "dropdown")
	}

	return elem.ListItem(
		vecty.Markup(
			classes,
		),
		elem.Anchor(
			vecty.Markup(
				prop.ID("packageDropdown"),
				prop.Href(""),
				vecty.Class("nav-link", "dropdown-toggle"),
				vecty.Property("role", "button"),
				vecty.Data("toggle", "dropdown"),
				vecty.Property("aria-haspopup", "true"),
				vecty.Property("aria-expanded", "false"),
				event.Click(func(ev *vecty.Event) {}).PreventDefault(),
			),
			vecty.Text(v.app.Scanner.DisplayName(v.app.Editor.CurrentPackage())),
		),
		elem.Div(
			packageItems...,
		),
	)
}
