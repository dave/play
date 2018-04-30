package stores

import (
	"errors"

	"bytes"
	"text/template"

	"strconv"

	"fmt"

	"github.com/dave/flux"
	"github.com/dave/play/actions"
	"github.com/dave/play/models"
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

func NewCompileStore(app *App) *CompileStore {
	s := &CompileStore{
		app: app,
	}
	return s
}

type CompileStore struct {
	app            *App
	compiling      bool
	compiled       bool
	consoleWritten bool
	tags           []string
}

func (s *CompileStore) Tags() []string {
	return s.tags
}

func (s *CompileStore) Compiling() bool {
	return s.compiling
}

func (s *CompileStore) Compiled() bool {
	return s.compiled
}

func (s *CompileStore) Handle(payload *flux.Payload) bool {
	switch payload.Action.(type) {
	case *actions.CompileStart:
		if err := s.compile(); err != nil {
			s.app.Fail(err)
			return true
		}
		payload.Notify()
	}
	return true
}

func (s *CompileStore) compile() error {
	path, count := s.app.Scanner.Main()
	if path == "" {
		if count == 0 {
			return errors.New("project has no main package")
		} else {
			return fmt.Errorf("project has %d main packages - select one and retry", count)
		}
	}

	if !s.app.Archive.Fresh(path) {
		s.app.Dispatch(
			&actions.RequestStart{Type: models.UpdateRequest, Run: true},
		)
		return nil
	}

	s.compiling = true
	defer func() {
		s.compiling = false
	}()

	s.app.Log("compiling")

	deps, err := s.app.Archive.Compile(path, s.Tags())
	if err != nil {
		return err
	}

	s.app.Log("running")

	doc := dom.GetWindow().Document()
	holder := doc.GetElementByID("iframe-holder")
	for _, v := range holder.ChildNodes() {
		v.Underlying().Call("remove")
	}
	frame := doc.CreateElement("iframe").(*dom.HTMLIFrameElement)
	frame.SetID("iframe")
	frame.Style().Set("width", "100%")
	frame.Style().Set("height", "100%")
	frame.Style().Set("border", "0")

	// We need to wait for the iframe to load before adding contents or Firefox will clear the iframe
	// after momentarily flashing up the contents.
	c := make(chan struct{})
	frame.AddEventListener("load", false, func(event dom.Event) {
		close(c)
	})

	holder.AppendChild(frame)
	<-c

	console := dom.GetWindow().Document().GetElementByID("console")
	console.SetInnerHTML("")
	frame.Get("contentWindow").Set("goPrintToConsole", js.InternalObject(func(b []byte) {
		console.SetInnerHTML(console.InnerHTML() + string(b))
		if !s.consoleWritten {
			s.consoleWritten = true
			s.app.Dispatch(&actions.ConsoleFirstWrite{})
		}
	}))

	frameDoc := frame.ContentDocument()

	if index, ok := s.app.Source.Files(path)["index.jsgo.html"]; ok {
		// has index

		indexTemplate, err := template.New("index").Parse(index)
		if err != nil {
			return err
		}
		data := struct{ Script string }{Script: ""}
		buf := &bytes.Buffer{}
		if err := indexTemplate.Execute(buf, data); err != nil {
			return err
		}

		frameDoc.Underlying().Call("open")
		frameDoc.Underlying().Call("write", buf.String())
		frameDoc.Underlying().Call("close")

		c := make(chan struct{})
		frame.AddEventListener("load", false, func(event dom.Event) {
			close(c)
		})
		<-c
	}

	head := frameDoc.GetElementsByTagName("head")[0].(*dom.BasicHTMLElement)

	loaderJs := ""
	for _, dep := range deps {
		loaderJs += "$load[" + strconv.Quote(dep.Path) + "]();\n"
	}
	scriptLoad := frameDoc.CreateElement("script")
	scriptLoad.SetID("loader")
	scriptLoad.SetInnerHTML(`
		var $load = {};
		var $count = 0;
		var $total = ` + fmt.Sprint(len(deps)) + `;
		var $finished = function() {
			` + loaderJs + `
			$mainPkg = $packages[` + strconv.Quote(path) + `];
			$synthesizeMethods();
			$packages["runtime"].$init();
			$go($mainPkg.$init, []);
			$flushConsole();
		};
		var $done = function() {
			$count++;
			if ($count == $total) {
				$finished();
			}
		};
	`)
	head.AppendChild(scriptLoad)

	for _, dep := range deps {
		scriptDep := frameDoc.CreateElement("script")
		scriptDep.SetID(dep.Path)
		scriptDep.SetInnerHTML(string(dep.Js) + "$done();")
		//scriptDep.AppendChild(doc.CreateTextNode(string(dep.Js) + "$done();"))
		head.AppendChild(scriptDep)
	}

	s.compiled = true
	s.app.Log()
	return nil
}
