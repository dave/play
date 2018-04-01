package actions

import (
	"github.com/dave/dropper"
	"github.com/dave/flux"
	"github.com/dave/jsgo/server/messages"
)

type Load struct{}

type ConsoleFirstWrite struct{}
type ConsoleToggleClick struct{}

type ChangeSplit struct{ Sizes []float64 }
type ChangeFile struct {
	Path string
	Name string
}

type LoadSource struct {
	Source         map[string]map[string]string
	CurrentPackage string
	CurrentFile    string
}

type UserChangedSplit struct{ Sizes []float64 }
type UserChangedText struct {
	Text    string
	Changed bool
}
type UserChangedFile struct{ Name string }
type UserChangedPackage struct{ Path string }

type AddFileClick struct{}
type AddPackageClick struct{}
type DeleteFileClick struct{}
type RemovePackageClick struct{}

type DownloadClick struct{}

type AddFile struct{ Name string }
type AddPackage struct{ Path string }
type DeleteFile struct{ Name string }
type RemovePackage struct{ Path string }

type FormatCode struct{ Then flux.ActionInterface }

// CompileStart compiles the app and injects the js into the iframe
type CompileStart struct{}

type DragEnter struct{}
type DragLeave struct{}
type DragDrop struct {
	Files   []dropper.File
	Changed map[string]map[string]bool
}

type Send struct{ Message messages.Message }
type Dial struct {
	Url     string
	Open    func() flux.ActionInterface
	Message func(interface{}) flux.ActionInterface
	Close   func() flux.ActionInterface
}

type ShareStart struct{}
type ShareOpen struct{}
type ShareMessage struct{ Message interface{} }
type ShareClose struct{}

type DeployStart struct{}
type DeployOpen struct{}
type DeployMessage struct{ Message interface{} }
type DeployClose struct{}

// UpdateStart updates the deps from the server and if Run == true, compiles and runs the app
type UpdateStart struct{ Run bool }
type UpdateOpen struct{ Main string }
type UpdateMessage struct{ Message interface{} }
type UpdateClose struct {
	Run  bool
	Main string
}

type GetStart struct{ Path string }
type GetOpen struct{ Path string }
type GetMessage struct {
	Path    string
	Message interface{}
}
type GetClose struct{}
