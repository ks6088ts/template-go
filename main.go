package main

import (
	"fmt"

	"github.com/ks6088ts/template-go/internal"
)

func main() {
	fmt.Printf("Hello, world.\nversion=%s, revision=%s\n", internal.Version, internal.Revision)
}
