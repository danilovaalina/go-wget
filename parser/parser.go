package main

import (
	"bytes"
	"io"

	"golang.org/x/net/html"
)

// ResourceURLs содержит все найденные URL-адреса
type ResourceURLs struct {
	Links   []string // <a href>
	Images  []string // <img src>
	Styles  []string // <link rel="stylesheet" href>
	Scripts []string // <script src>
}

// GetAllResources returns all URLs in a single slice (useful for deduplication and processing)
func (r *ResourceURLs) GetAllResources() []string {
	all := make([]string, 0, len(r.Links)+len(r.Images)+len(r.Styles)+len(r.Scripts))
	all = append(all, r.Links...)
	all = append(all, r.Images...)
	all = append(all, r.Styles...)
	all = append(all, r.Scripts...)
	return all
}

func ParseHTML(htmlData []byte) (*ResourceURLs, error) {
	tokenizer := html.NewTokenizer(bytes.NewReader(htmlData))

	var urls ResourceURLs

	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			err := tokenizer.Err()
			if err == io.EOF {
				return &urls, nil
			}
			return nil, err

		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			tag := token.Data

			// Извлекаем атрибуты в map для удобства
			attrs := make(map[string]string)
			for _, attr := range token.Attr {
				attrs[attr.Key] = attr.Val
			}

			switch tag {
			case "a":
				if href, ok := attrs["href"]; ok && href != "" {
					urls.Links = append(urls.Links, href)
				}

			case "img":
				if src, ok := attrs["src"]; ok && src != "" {
					urls.Images = append(urls.Images, src)
				}

			case "link":
				// Только CSS: rel="stylesheet"
				if rel, ok := attrs["rel"]; ok && rel == "stylesheet" {
					if href, ok := attrs["href"]; ok && href != "" {
						urls.Styles = append(urls.Styles, href)
					}
				}

			case "script":
				if src, ok := attrs["src"]; ok && src != "" {
					urls.Scripts = append(urls.Scripts, src)
				}
			}
		default:
			continue
		}
	}
}
