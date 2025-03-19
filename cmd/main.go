package main

import (
	"bookie/pdf"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	output := flag.String("o", "", "Path to the fb2 book")
	skip := flag.Bool("skip-unknown", true, "Skip unknown XML elements (may contain spam)")

	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Printf("Error: Expected exactly one book path argument, got %d\n", len(args))
		fmt.Printf("Usage: %s [options] <path_to_fb2>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	var outPath string
	if *output != "" {
		outPath = *output
	} else {
		if strings.HasSuffix(args[0], ".fb2") {
			outPath = strings.TrimSuffix(args[0], "fb2") + "pdf"
		} else {
			outPath += args[0] + ".pdf"
		}
	}

	err := pdf.NewConverter(*skip).WritePDF(args[0], outPath)
	if err != nil {
		log.Fatal(err)
	}
}
