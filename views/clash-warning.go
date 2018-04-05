package views

import (
	"fmt"
	"sort"

	"github.com/dave/play/models"
	"github.com/dave/play/stores"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

type ClashWarningModal struct {
	*Modal
}

func NewClashWarningModal(app *stores.App) *ClashWarningModal {
	v := &ClashWarningModal{
		&Modal{
			app:    app,
			id:     models.ClashWarningModal,
			title:  "Warning",
			action: nil,
		},
	}
	return v
}

func (v *ClashWarningModal) Render() vecty.ComponentOrHTML {
	var paragraphs []vecty.MarkupOrChild
	for source, compiled := range v.app.Scanner.Clashes() {
		var paths []string
		for p := range compiled {
			paths = append(paths, p)
		}
		var text string
		if len(compiled) == 0 {
			continue
		} else {
			sort.Strings(paths)
			if len(compiled) == 1 {
				text = fmt.Sprintf("Source package %s is imported by this pre-compiled package:", source)
			} else {
				text = fmt.Sprintf("Source package %s is imported by these pre-compiled packages:", source)
			}
			var items []vecty.MarkupOrChild
			for _, path := range paths {
				items = append(items, elem.ListItem(
					vecty.Text(path),
				))
			}
			paragraphs = append(paragraphs, elem.Paragraph(
				vecty.Text(text),
				elem.UnorderedList(items...),
			))
		}

	}
	var text string
	if len(v.app.Scanner.Clashes()) == 1 {
		text = "If the external interface this source packages changes, your code may break at run-time. To solve this, either load the source for all pre-compiled packages, or use the Update feature each time you change the external interface of the source package."
	} else {
		text = "If the external interface these source packages change, your code may break at run-time. To solve this, either load the source for all pre-compiled packages, or use the Update feature each time you change the external interface of the source packages."
	}
	paragraphs = append(paragraphs, elem.Paragraph(
		vecty.Text(text),
	))

	return v.Body(
		paragraphs...,
	).Build()
}
