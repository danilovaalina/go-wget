package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
)

func main() {
	var depth int
	var output string

	flag.IntVar(&depth, "depth", 1, "Maximum recursion depth")
	flag.StringVar(&output, "output", "mirror", "Output directory")
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatal("Usage: wget <URL> [--depth N] [--output DIR]")
	}

	startURL := flag.Arg(0)
	parsedURL, err := url.Parse(startURL)
	if err != nil {
		log.Fatalf("Invalid URL: %v", err)
	}

	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "https"
	}

	fmt.Printf("Starting download of %s\n", parsedURL.String())
	fmt.Printf("Depth: %d, Output: %s\n", depth, output)
}
