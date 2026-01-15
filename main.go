package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"
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

	// После парсинга startURL
	downloader := NewDownloader(DefaultOptions)
	result, err := downloader.Fetch(parsedURL)
	if err != nil {
		log.Fatalf("Download failed: %v", err)
	}
	fmt.Printf("Downloaded %d bytes from %s\n", len(result.Body), result.URL)

	// После получения result от downloader
	saver, err := NewSaver(output)
	if err != nil {
		log.Fatalf("Failed to create saver: %v", err)
	}

	savedPath, err := saver.Save(result.URL, result.Body)
	if err != nil {
		log.Fatalf("Failed to save: %v", err)
	}

	fmt.Printf("Saved to: %s\n", savedPath)

	if strings.Contains(string(result.Body), "<html") {
		parsed, err := ParseHTML(result.Body)
		if err != nil {
			log.Printf("Failed to parse HTML: %v", err)
		} else {
			fmt.Printf("Found %d links, %d images, %d styles, %d scripts\n",
				len(parsed.Links), len(parsed.Images), len(parsed.Styles), len(parsed.Scripts))
			for _, link := range parsed.Links[:min(5, len(parsed.Links))] {
				fmt.Printf(" → Link: %s\n", link)
			}
		}
	}
}
