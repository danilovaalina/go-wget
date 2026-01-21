// cmd/wget/main.go
package main

import (
	"flag"
	"log"
	"net/url"

	"go-wget/crawler"
	"go-wget/downloader"
	"go-wget/saver"
)

// Config содержит параметры запуска утилиты
type Config struct {
	URL    *url.URL
	Depth  int
	Output string
}

// parseFlags парсит аргументы командной строки и возвращает конфигурацию
func parseFlags() *Config {
	var depth int
	var output string

	flag.IntVar(&depth, "depth", 1, "Maximum recursion depth")
	flag.StringVar(&output, "output", "mirror", "Output directory")
	flag.Usage = func() {
		log.Print("Usage: wget [flags] <URL>")
		log.Print("Example: wget --depth 2 https://example.com")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		log.Fatal("error: exactly one URL must be provided")
	}

	startURL := flag.Arg(0)
	parsedURL, err := url.Parse(startURL)
	if err != nil {
		log.Fatalf("error: invalid URL %q: %v", startURL, err)
	}

	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "https"
	}

	return &Config{
		URL:    parsedURL,
		Depth:  depth,
		Output: output,
	}
}

func main() {
	cfg := parseFlags()

	dl := downloader.New(downloader.DefaultOptions)
	sv, err := saver.New(cfg.Output)
	if err != nil {
		log.Fatalf("Failed to create saver: %v", err)
	}

	c := crawler.New(dl, sv, cfg.URL)
	if err := c.Crawl(cfg.URL, cfg.Depth); err != nil {
		log.Fatalf("Crawling failed: %v", err)
	}

	log.Println("Mirroring completed.")
}
