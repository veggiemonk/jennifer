package main

import (
	"bytes"
	"flag"
	"log"
	"os"
)

func main() {
	// notest

	out := flag.String("output", "./package-names.go", "Output filename to write")
	pkg := flag.String("package", "main", "Package name in generated file")
	name := flag.String("name", "PackageNames", "Name of the variable to define")
	filter := flag.String("filter", ".*", "Regex to filter paths (operates on full path including vendor directory)")
	standard := flag.Bool("standard", false, "Use standard library packages")
	novendor := flag.Bool("novendor", false, "Exclude packages in vendor directories")
	goListPath := flag.String("path", "all", "Path to pass to go list command")
	flag.Parse()

	buf := &bytes.Buffer{}
	if err := hints(buf, *pkg, *name, *goListPath, *filter, *standard, *novendor); err != nil {
		log.Fatal(err.Error())
	}
	if err := os.WriteFile(*out, buf.Bytes(), 0o644); err != nil {
		log.Fatal(err.Error())
	}
}
