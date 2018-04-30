package views

import (
	"github.com/dave/play/actions"
	"github.com/dave/play/models"
	"github.com/dave/play/stores"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/russross/blackfriday"
)

type HelpModal struct {
	*Modal
}

func NewHelpModal(app *stores.App) *HelpModal {
	v := &HelpModal{}
	v.Modal = &Modal{
		app:    app,
		id:     models.HelpModal,
		title:  "Help",
		action: v.action,
		large:  true,
	}
	return v
}

func (v *HelpModal) Render() vecty.ComponentOrHTML {

	// Render the markdown input into HTML using Blackfriday.
	unsafeHTML := blackfriday.MarkdownCommon([]byte(helpMarkdown))

	// Sanitize the HTML.
	//safeHTML := string(bluemonday.UGCPolicy().SanitizeBytes(unsafeHTML))

	return v.Body(
		elem.Div(
			vecty.Markup(
				vecty.UnsafeHTML(string(unsafeHTML)),
			),
		),
	).Build()
}

func (v *HelpModal) action(*vecty.Event) {
	v.app.Dispatch(&actions.ModalClose{Modal: models.HelpModal})
}

var helpMarkdown = `
Edit and run Go code, compiled to JS using GopherJS, supporting arbitrary import paths!

https://play.jsgo.io/

[<img width="500" alt="title936803092" src="https://user-images.githubusercontent.com/925351/39423295-7b461464-4c71-11e8-9565-6c860e8642e8.png">](https://play.jsgo.io/)

## Demos
* https://play.jsgo.io/github.com/hajimehoshi/ebiten/examples/flappy
* https://play.jsgo.io/github.com/hajimehoshi/ebiten/examples/2048
* https://play.jsgo.io/github.com/shurcooL/tictactoe/cmd/tictactoe
* https://play.jsgo.io/github.com/dave/compress/zip
* https://play.jsgo.io/github.com/dave/jstest

## Features

### Initialise
The URL can be used to initialise with code in several ways:

* Load a Go package with ` + "`" + `/{{ Package path }}` + "`" + `
* Load a Github Gist with ` + "`" + `/gist.github.com/{{ Gist ID }}` + "`" + `
* Load a shared project with ` + "`" + `/{{ Share ID }}` + "`" + `
* Load a ` + "`" + `play.golang.org` + "`" + ` share with ` + "`" + `/p/{{ Go playground ID }}` + "`" + `

### Run
<img width="150" alt="run" src="https://user-images.githubusercontent.com/925351/39422110-550c650a-4c6c-11e8-9353-050f823c6201.png">

Click the ` + "`" + `Run` + "`" + ` button to run your code in the right-hand panel. If the imports have been changed recently,
the dependencies will be refreshed before running.

### Format code
<img width="150" alt="format" src="https://user-images.githubusercontent.com/925351/39422105-54677d7e-4c6c-11e8-8cfa-3b7013d6cf64.png">

Use the ` + "`" + `Format code` + "`" + ` option to run ` + "`" + `gofmt` + "`" + ` on your code. This is executed automatically when the ` + "`" + `Run` + "`" + `, 
` + "`" + `Update` + "`" + `, ` + "`" + `Share` + "`" + ` or ` + "`" + `Deploy` + "`" + ` features are used.

### Update
<img width="150" alt="update" src="https://user-images.githubusercontent.com/925351/39422115-557afea2-4c6c-11e8-9af5-fb98f582ae6d.png">

If you update a dependency, use the ` + "`" + `Update` + "`" + ` option, which does the equivalent of ` + "`" + `go get -u` + "`" + ` and refreshes 
the changes in any import or dependency.   

### Share
<img width="150" alt="share" src="https://user-images.githubusercontent.com/925351/39422111-55268a34-4c6c-11e8-8ab6-bb6f718bcf5d.png">

To share your project with others, use the ` + "`" + `Share` + "`" + ` option. Your project will be persisted to a json file 
on ` + "`" + `src.jsgo.io` + "`" + ` and the page will update to a sharable URL.

### Deploy
<img width="150" alt="deploy" src="https://user-images.githubusercontent.com/925351/39422100-53ddecf8-4c6c-11e8-820c-4115472d4b8c.png"><img width="150" alt="deploy2-250" src="https://user-images.githubusercontent.com/925351/39422101-53f5fb9a-4c6c-11e8-9f80-90914eab4dd1.png">

To deploy your code to [jsgo.io](https://jsgo.io), use the ` + "`" + `Deploy` + "`" + ` feature. A modal will be displayed with the 
link to the page on ` + "`" + `jsgo.io` + "`" + `, and the Loader JS on ` + "`" + `pkg.jsgo.io` + "`" + `. 

Use the ` + "`" + `jsgo.io` + "`" + ` link for testing and toy projects. Remember you're sharing the ` + "`" + `jsgo.io` + "`" + ` domain with 
everyone else, so the browser environment should be considered toxic.

The Loader JS on ` + "`" + `pkg.jsgo.io` + "`" + ` can be used in production, and should be added to a script tag on your 
own website. See [github.com/dave/jsgo](https://github.com/dave/jsgo) for more information.

### Console
<img width="150" alt="console" src="https://user-images.githubusercontent.com/925351/39422096-53904c3c-4c6c-11e8-94f6-2c8f62c1f9a3.png"><img width="150" alt="console1" src="https://user-images.githubusercontent.com/925351/39422097-53ac0026-4c6c-11e8-9cbf-aa08411ff02d.png">

Writes to ` + "`" + `os.Stdout` + "`" + ` are redirected to a playground console, which can be toggled using the ` + "`" + `Show console` + "`" + `
option. The console will automatically appear the first time it's written to.

### Minify
<img width="150" alt="minify" src="https://user-images.githubusercontent.com/925351/39422107-54a89b56-4c6c-11e8-9eba-8d6c5492fef3.png">

In normal usage, all JS is minified. For debugging, this can be toggled with the ` + "`" + `Minify JS` + "`" + ` option.

### Build tags
<img width="150" alt="tags" src="https://user-images.githubusercontent.com/925351/39422112-554002fc-4c6c-11e8-8c2a-79b8e1f13045.png"><img width="150" alt="tags2" src="https://user-images.githubusercontent.com/925351/39422114-5560a214-4c6c-11e8-987b-842620eb7abe.png">

The build tags used when compiling can be edited with the ` + "`" + `Build tags...` + "`" + ` option. The selected build 
tags are persisted when using the ` + "`" + `Share` + "`" + ` feature.

### Download
<img width="150" alt="download" src="https://user-images.githubusercontent.com/925351/39422103-54358530-4c6c-11e8-8dbb-23b109bab9f8.png">

The ` + "`" + `Download` + "`" + ` option downloads the project. Single file projects are downloaded as a single file, while
multi-file projects download as a zip.

### Upload
<img width="150" alt="download" src="https://user-images.githubusercontent.com/925351/39422494-e021814c-4c6d-11e8-8dd4-3ceb330d6d97.png">

Files can be uploaded to the project simply by drag+drop. Zip files generated by the ` + "`" + `Download` + "`" + ` feature
can be uploaded to restore a multi-file project.

### File menu
<img width="150" alt="files" src="https://user-images.githubusercontent.com/925351/39422104-544e3f8a-4c6c-11e8-9953-002ae51db341.png">

Change the selected file with the file menu.

### Add file
<img width="150" alt="add-file" src="https://user-images.githubusercontent.com/925351/39422092-535a7d46-4c6c-11e8-9634-b9c36bb7b943.png">

Add a file to the current package with the ` + "`" + `Add file` + "`" + ` option. Only ` + "`" + `.go` + "`" + `, ` + "`" + `.md` + "`" + ` and ` + "`" + `.inc.js` + "`" + ` files are
supported. If no extension is supplied, ` + "`" + `.go` + "`" + ` is added.

### Delete file
<img width="150" alt="delete-file" src="https://user-images.githubusercontent.com/925351/39422099-53c48e48-4c6c-11e8-838a-4f07a7db41bf.png">

Delete a file from the current package with the ` + "`" + `Delete file` + "`" + ` option.

### Package menu
<img width="150" alt="package" src="https://user-images.githubusercontent.com/925351/39422108-54cd5946-4c6c-11e8-9de8-21c9ff7a8bb3.png">

Change the selected package with the package menu.

### Add package
<img width="150" alt="add-package" src="https://user-images.githubusercontent.com/925351/39422094-53742566-4c6c-11e8-86ad-9f2753c12c33.png">

Add an empty package with the ` + "`" + `Add package` + "`" + ` option.

### Load package
<img width="150" alt="load-package" src="https://user-images.githubusercontent.com/925351/39422106-5480ef70-4c6c-11e8-97ce-cb169cf2219b.png"><img width="150" alt="load-package-1" src="https://user-images.githubusercontent.com/925351/39422297-0fdaff90-4c6d-11e8-9d72-0f21e07e5d3a.png"><img width="150" alt="load-package-2" src="https://user-images.githubusercontent.com/925351/39422298-0ffaa91c-4c6d-11e8-8a09-bb645907c503.png">

The source for an import or dependency can be loaded with the ` + "`" + `Load package` + "`" + ` option. By default, only 
the direct imports of your project are listed. Use the ` + "`" + `Show all dependencies` + "`" + ` option to show the entire
dependency tree.

### Remove package
<img width="150" alt="remove-package" src="https://user-images.githubusercontent.com/925351/39422109-54e9836e-4c6c-11e8-913b-425945c194e0.png">

A package can be removed with the ` + "`" + `Remove package` + "`" + ` option.

## How to contact me

If you'd like to chat more about the project, feel free to [add an issue](https://github.com/dave/play/issues), 
mention [@dave](https://github.com/dave/) in your PR, email me or post in the #gopherjs channel of the 
Gophers Slack. I'm happy to help!
`
