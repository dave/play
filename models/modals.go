package models

type Modal string

const (
	ClashWarningModal  Modal = "clash-warning-modal"
	AddPackageModal    Modal = "add-package-modal"
	AddFileModal       Modal = "add-file-modal"
	DeleteFileModal    Modal = "delete-file-modal"
	DeployDoneModal    Modal = "deploy-done-modal"
	LoadPackageModal   Modal = "load-package-modal"
	RemovePackageModal Modal = "remove-package-modal"
	BuildTagsModal     Modal = "build-tags-modal"
	HelpModal          Modal = "help-modal"
)

type RequestType string

const (
	GetRequest        RequestType = "get"
	UpdateRequest     RequestType = "update"
	InitialiseRequest RequestType = "initialise"
)
