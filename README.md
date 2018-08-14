<a href="https://patreon.com/davebrophy" title="Help with my hosting bills using Patreon"><img src="https://img.shields.io/badge/patreon-donate-yellow.svg" style="max-width:100%;"></a>

Edit and run Go in the browser, supporting arbitrary import paths!

https://play.jsgo.io/

[<img width="500" alt="title936803092" src="https://user-images.githubusercontent.com/925351/39423295-7b461464-4c71-11e8-9565-6c860e8642e8.png">](https://play.jsgo.io/)

The jsgo playground is an extension of the jsgo compiler. The compiler allows you to easily compile Go 
to JS using GopherJS, and automatically host the results in an aggressively cached CDN. The playground 
adds an online editor and many other features (see below).

The unique feature of the jsgo playground is that it supports arbitrary import paths. Other Go playgrounds 
are limited to just the Go standard library.

For more for more info:

* jsgo compiler: https://github.com/dave/jsgo  
* jsgo playground: https://github.com/dave/play  

## Demos

Here's the simplest demo - it just writes to the console and to the page:  

* https://play.jsgo.io/github.com/dave/jstest

Here's a couple of simple demos that accept files by drag and drop. The first compresses dropped files to 
a zip. The second compresses images to jpg. They use the Go standard library zip / image libraries, which 
work flawlessly in the browser:   

* https://play.jsgo.io/github.com/dave/zip
* https://play.jsgo.io/github.com/dave/img 

The amazing ebiten 2D games library is a perfect example of the power of Go in the browser. Here's some 
demos: 

* https://play.jsgo.io/github.com/hajimehoshi/ebiten/examples/2048
* https://play.jsgo.io/github.com/hajimehoshi/go-inovation
* https://play.jsgo.io/github.com/hajimehoshi/ebiten/examples/flappy

## Contact

If you'd like to chat more about the project, feel free to [add an issue](https://github.com/dave/play/issues), 
mention [@dave](https://github.com/dave/) or post in the #gopherjs channel of the Gophers Slack. I'm 
happy to help!

## Features

#### Initialise
The URL can be used to initialise with code in several ways:

* Load a Go package with `/{{ Package path }}`
* Load a Github Gist with `/gist.github.com/{{ Gist ID }}`
* Load a shared project with `/{{ Share ID }}`
* Load a `play.golang.org` share with `/p/{{ Go playground ID }}`

<img align="right" width="150" alt="run" src="https://user-images.githubusercontent.com/925351/39422110-550c650a-4c6c-11e8-9353-050f823c6201.png">

#### Run
Click the `Run` button to run your code in the right-hand panel. If the imports have been changed recently,
the dependencies will be refreshed before running.

<table></table>

<img align="right" width="150" alt="format" src="https://user-images.githubusercontent.com/925351/39422105-54677d7e-4c6c-11e8-8cfa-3b7013d6cf64.png">

#### Format code
Use the `Format code` option to run `gofmt` on your code. This is executed automatically when the `Run`, 
`Update`, `Share` or `Deploy` features are used.

<table></table>

<img align="right" width="150" alt="update" src="https://user-images.githubusercontent.com/925351/39422115-557afea2-4c6c-11e8-9af5-fb98f582ae6d.png">

#### Update
If you update a dependency, use the `Update` option, which does the equivalent of `go get -u` and refreshes 
the changes in any import or dependency.   

<table></table>

<img align="right" width="150" alt="share" src="https://user-images.githubusercontent.com/925351/39422111-55268a34-4c6c-11e8-8ab6-bb6f718bcf5d.png">

#### Share
To share your project with others, use the `Share` option. Your project will be persisted to a json file 
on `src.jsgo.io` and the page will update to a sharable URL.

<table></table>

<img align="right" width="150" alt="deploy" src="https://user-images.githubusercontent.com/925351/39422100-53ddecf8-4c6c-11e8-820c-4115472d4b8c.png">

#### Deploy
To deploy your code to [jsgo.io](https://jsgo.io), use the `Deploy` feature. A modal will be displayed with the 
link to the page on `jsgo.io`, and the Loader JS on `pkg.jsgo.io`. 

Use the `jsgo.io` link for testing and toy projects. Remember you're sharing the `jsgo.io` domain with 
everyone else, so the browser environment should be considered toxic.

The Loader JS on `pkg.jsgo.io` can be used in production, and should be added to a script tag on your 
own website. See [github.com/dave/jsgo](https://github.com/dave/jsgo) for more information.

<table></table>

<img align="right" width="150" alt="console" src="https://user-images.githubusercontent.com/925351/39422096-53904c3c-4c6c-11e8-94f6-2c8f62c1f9a3.png">

#### Console
Writes to `os.Stdout` are redirected to a playground console, which can be toggled using the `Show console`
option. The console will automatically appear the first time it's written to.

<table></table>

<img align="right" width="150" alt="minify" src="https://user-images.githubusercontent.com/925351/39422107-54a89b56-4c6c-11e8-9eba-8d6c5492fef3.png">

#### Minify
In normal usage, all JS is minified. For debugging, this can be toggled with the `Minify JS` option.

<table></table>

<img align="right" width="150" alt="tags" src="https://user-images.githubusercontent.com/925351/39422112-554002fc-4c6c-11e8-8c2a-79b8e1f13045.png">

#### Build tags
The build tags used when compiling can be edited with the `Build tags...` option. The selected build 
tags are persisted when using the `Share` feature.

<table></table>

<img align="right" width="150" alt="download" src="https://user-images.githubusercontent.com/925351/39422103-54358530-4c6c-11e8-8dbb-23b109bab9f8.png">

#### Download
The `Download` option downloads the project. Single file projects are downloaded as a single file, while
multi-file projects download as a zip.

<table></table>

<img align="right" width="150" alt="download" src="https://user-images.githubusercontent.com/925351/39422494-e021814c-4c6d-11e8-8dd4-3ceb330d6d97.png">

#### Upload
Files can be uploaded to the project simply by drag+drop. Zip files generated by the `Download` feature
can be uploaded to restore a multi-file project.

<table></table>

<img align="right" width="150" alt="files" src="https://user-images.githubusercontent.com/925351/39422104-544e3f8a-4c6c-11e8-9953-002ae51db341.png">

#### File menu
Change the selected file with the file menu.

<table></table>

<img align="right" width="150" alt="add-file" src="https://user-images.githubusercontent.com/925351/39422092-535a7d46-4c6c-11e8-9634-b9c36bb7b943.png">

#### Add file
Add a file to the current package with the `Add file` option. Only `.go`, `.md` and `.inc.js` files are
supported. If no extension is supplied, `.go` is added.

<table></table>

<img align="right" width="150" alt="delete-file" src="https://user-images.githubusercontent.com/925351/39422099-53c48e48-4c6c-11e8-838a-4f07a7db41bf.png">

#### Delete file
Delete a file from the current package with the `Delete file` option.

<table></table>

<img align="right" width="150" alt="package" src="https://user-images.githubusercontent.com/925351/39422108-54cd5946-4c6c-11e8-9de8-21c9ff7a8bb3.png">

#### Package menu
Change the selected package with the package menu.

<table></table>

<img align="right" width="150" alt="add-package" src="https://user-images.githubusercontent.com/925351/39422094-53742566-4c6c-11e8-86ad-9f2753c12c33.png">

#### Add package
Add an empty package with the `Add package` option.

<table></table>

<img align="right" width="150" alt="load-package" src="https://user-images.githubusercontent.com/925351/39422106-5480ef70-4c6c-11e8-97ce-cb169cf2219b.png">

#### Load package
The source for an import or dependency can be loaded with the `Load package` option. By default, only 
the direct imports of your project are listed. Use the `Show all dependencies` option to show the entire
dependency tree.

<table></table>

<img align="right" width="150" alt="remove-package" src="https://user-images.githubusercontent.com/925351/39422109-54e9836e-4c6c-11e8-913b-425945c194e0.png">

#### Remove package
A package can be removed with the `Remove package` option.

