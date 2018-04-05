package main

import (
	"fmt"

	"github.com/dave/play/testing/a"
	"github.com/dave/play/testing/b"
)

func main() {
	fmt.Println(b.B(a.A{A: 1}))
}
