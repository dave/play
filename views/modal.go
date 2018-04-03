package views

import (
	"github.com/dave/play/actions"
	"github.com/dave/play/models"
	"github.com/dave/play/stores"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
)

type Modal struct {
	vecty.Core
	app    *stores.App
	id     models.Modal
	title  string
	action func(*vecty.Event)
	body   []vecty.MarkupOrChild // This should be set on every Render
	shown  func()
	hidden func()
}

func (m *Modal) Mount() {
	m.app.Watch(m, func(done chan struct{}) {
		defer close(done)
		modalIsVisible := js.Global.Call("$", "#"+m.id).Call("hasClass", "show").Bool()
		shouldBeVisible := m.app.Page.ModalOpen(m.id)
		if modalIsVisible != shouldBeVisible {
			if shouldBeVisible {
				js.Global.Call("$", "#"+m.id).Call("modal", "show")
			} else {
				js.Global.Call("$", "#"+m.id).Call("modal", "hide")
			}
		}
	})
	js.Global.Call("$", "#"+m.id).Call("on", "shown.bs.modal", func() {
		if !m.app.Page.ModalOpen(m.id) {
			m.app.Dispatch(&actions.ModalOpen{Modal: m.id})
		}
		if m.shown != nil {
			m.shown()
		}
	})
	js.Global.Call("$", "#"+m.id).Call("on", "hidden.bs.modal", func() {
		if m.app.Page.ModalOpen(m.id) {
			m.app.Dispatch(&actions.ModalClose{Modal: m.id})
		}
		if m.hidden != nil {
			m.hidden()
		}
	})
}

func (m *Modal) Unmount() {
	m.app.Delete(m)
	js.Global.Call("$", "#"+m.id).Call("unbind")
}

func (m *Modal) Body(body ...vecty.MarkupOrChild) *Modal {
	m.body = body
	return m
}

func (m *Modal) Build() vecty.ComponentOrHTML {

	body := []vecty.MarkupOrChild{
		vecty.Markup(
			vecty.Class("modal-body"),
		),
	}
	body = append(body, m.body...)

	okDisplay := ""
	if m.action == nil {
		okDisplay = "none"
	}

	return elem.Div(
		vecty.Markup(
			prop.ID(string(m.id)),
			vecty.Class("modal"),
			vecty.Property("tabindex", "-1"),
			vecty.Property("role", "dialog"),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("modal-dialog"),
				vecty.Property("role", "dialog"),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("modal-content"),
				),
				elem.Div(
					vecty.Markup(
						vecty.Class("modal-header"),
					),
					elem.Heading5(
						vecty.Markup(
							vecty.Class("modal-title"),
						),
						vecty.Text(m.title),
					),
					elem.Button(
						vecty.Markup(
							prop.Type(prop.TypeButton),
							vecty.Class("close"),
							vecty.Data("dismiss", "modal"),
							vecty.Property("aria-label", "Close"),
						),
						elem.Span(
							vecty.Markup(
								vecty.Property("aria-hidden", "true"),
							),
							vecty.Text("×"),
						),
					),
				),
				elem.Div(
					body...,
				),
				elem.Div(
					vecty.Markup(
						vecty.Class("modal-footer"),
					),
					elem.Button(
						vecty.Markup(
							prop.Type(prop.TypeButton),
							vecty.Class("btn", "btn-primary"),
							event.Click(m.action).PreventDefault(),
							vecty.Style("display", okDisplay),
						),
						vecty.Text("OK"),
					),
					elem.Button(
						vecty.Markup(
							prop.Type(prop.TypeButton),
							vecty.Class("btn", "btn-secondary"),
							vecty.Data("dismiss", "modal"),
						),
						vecty.Text("Close"),
					),
				),
			),
		),
	)
}
