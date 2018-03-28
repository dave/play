package views

import (
	"github.com/dave/play/actions"
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

	var fileItems []vecty.MarkupOrChild
	fileItems = append(fileItems,
		vecty.Markup(
			vecty.Class("dropdown-menu"),
			vecty.Property("aria-labelledby", "fileDropdown"),
		),
	)
	for _, name := range v.app.Editor.Filenames() {
		name := name
		fileItems = append(fileItems,
			elem.Anchor(
				vecty.Markup(
					vecty.Class("dropdown-item"),
					vecty.ClassMap{
						"disabled": name == v.app.Editor.Current(),
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
					v.app.Dispatch(&actions.AddFileClick{})
				}).PreventDefault(),
			),
			vecty.Text("Add file..."),
		),
		elem.Anchor(
			vecty.Markup(
				vecty.Class("dropdown-item"),
				prop.Href(""),
				event.Click(func(e *vecty.Event) {
					v.app.Dispatch(&actions.DeleteFileClick{})
				}).PreventDefault(),
			),
			vecty.Text("Delete file..."),
		),
	)
	fileDropdownClasses := vecty.Class("nav-item", "dropdown")
	if v.app.Editor.Current() == "" || len(v.app.Editor.Files()) <= 1 {
		fileDropdownClasses = vecty.Class("nav-item", "dropdown", "d-none")
	}

	return elem.Navigation(
		vecty.Markup(
			vecty.Class("menu", "navbar", "navbar-expand", "navbar-light", "bg-light"),
		),
		elem.UnorderedList(
			vecty.Markup(
				vecty.Class("navbar-nav", "mr-auto"),
			),
			/*
				elem.ListItem(
					vecty.Markup(
						vecty.Class("nav-item", "dropdown"),
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
						vecty.Text("main"),
					),
					elem.Div(
						vecty.Markup(
							vecty.Class("dropdown-menu"),
							vecty.Property("aria-labelledby", "packageDropdown"),
						),
						elem.Anchor(
							vecty.Markup(
								vecty.Class("dropdown-item"),
								prop.Href(""),
								event.Click(func(e *vecty.Event) {}).PreventDefault(),
							),
							vecty.Text("Package 1"),
						),
						elem.Anchor(
							vecty.Markup(
								vecty.Class("dropdown-item"),
								prop.Href(""),
								event.Click(func(e *vecty.Event) {}).PreventDefault(),
							),
							vecty.Text("Package 2"),
						),
						elem.Anchor(
							vecty.Markup(
								vecty.Class("dropdown-item"),
								prop.Href(""),
								event.Click(func(e *vecty.Event) {}).PreventDefault(),
							),
							vecty.Text("Package 3"),
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
								event.Click(func(e *vecty.Event) {}).PreventDefault(),
							),
							vecty.Text("Add package..."),
						),
					),
				),
			*/
			elem.ListItem(
				vecty.Markup(
					fileDropdownClasses,
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
					vecty.Text(v.app.Editor.Current()),
				),
				elem.Div(
					fileItems...,
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
						vecty.Style("margin-right", "10px"),
					),
					vecty.Text(""),
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
							} else if v.app.Archive.Fresh() {
								v.app.Dispatch(&actions.FormatCode{
									Then: &actions.CompileStart{},
								})
							} else {
								v.app.Dispatch(&actions.FormatCode{
									Then: &actions.UpdateStart{Run: true},
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
								v.app.Dispatch(&actions.FormatCode{
									Then: &actions.UpdateStart{},
								})
							}).PreventDefault(),
						),
						vecty.Text("Update"),
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
									Then: &actions.ShareStart{},
								})
							}).PreventDefault(),
						),
						vecty.Text("Share"),
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
								v.app.Dispatch(&actions.AddFileClick{})
							}).PreventDefault(),
						),
						vecty.Text("Add file..."),
					),
					/*

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
									js.Global.Call("alert", "TODO")
								}).PreventDefault(),
							),
							vecty.Text("Build tags..."),
						),
						elem.Anchor(
							vecty.Markup(
								vecty.Class("dropdown-item"),
								prop.Href(""),
								event.Click(func(e *vecty.Event) {
									js.Global.Call("alert", "TODO")
								}).PreventDefault(),
							),
							vecty.Text("Save"),
						),
						elem.Anchor(
							vecty.Markup(
								vecty.Class("dropdown-item"),
								prop.Href(""),
								event.Click(func(e *vecty.Event) {
									js.Global.Call("alert", "TODO")
								}).PreventDefault(),
							),
							vecty.Text("Deploy"),
						),*/
				),
			),
		),
	)
}
