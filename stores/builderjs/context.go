package builderjs

import (
	"go/build"
	"os"

	"io"

	"bytes"

	"io/ioutil"
	"path/filepath"
)

func NewBuildContext(source map[string]string, tags []string) *build.Context {

	tags = append(tags, "netgo", "purego", "jsgo")

	b := &build.Context{
		GOARCH:        "js",     // Target architecture
		GOOS:          "darwin", // Target operating system
		GOROOT:        "goroot", // Go root
		GOPATH:        "gopath", // Go path
		InstallSuffix: "",       // Builder only: "min" or "".
		Compiler:      "gc",     // Compiler to assume when computing target paths
		BuildTags:     tags,     // Build tags
		CgoEnabled:    false,    // Builder only: detect `import "C"` to throw proper error
		ReleaseTags:   build.Default.ReleaseTags,

		IsDir:     func(path string) bool { panic("should not be called in JS") },
		HasSubdir: func(root, dir string) (rel string, ok bool) { panic("should not be called in JS") },
		ReadDir:   func(path string) ([]os.FileInfo, error) { panic("should not be called in JS") },

		// OpenFile opens a file (not a directory) for reading.
		// If OpenFile is nil, Import uses os.Open.
		OpenFile: func(path string) (io.ReadCloser, error) {
			_, name := filepath.Split(path)
			s, ok := source[name]
			if !ok {
				return nil, os.ErrNotExist
			}
			return ioutil.NopCloser(bytes.NewBuffer([]byte(s))), nil
		},
	}
	return b
}
