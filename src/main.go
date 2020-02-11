package main

import (
	"flag"
	"fmt"
	"os"

	kitten "github.com/riavalon/gottabkitten"
)

func main() {
	sourceFile := flag.String("input", "", "Specify an input file with text to read from. If not supplied, expects sourceText")
	sourceText := flag.String("text", "", "Text to be used. If not specified, sourceFile is expected. SourceFile takes precendence.")
	flag.Parse()

	if len(*sourceFile) == 0 && len(*sourceText) == 0 {
		fmt.Println("You must specify either a source file or source text")
		os.Exit(1)
	}

	var textContent string
	if len(*sourceFile) > 0 {
		textContent = string(kitten.ReadSourceFile(*sourceFile))
	} else {
		textContent = *sourceText
	}

	kitten.Impurrove(textContent)
}
