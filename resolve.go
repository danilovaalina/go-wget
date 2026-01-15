package main

import (
	"fmt"
	"net/url"
)

// ResolveReference преобразует относительный URL в абсолютный,
// используя base как контекст (как это делает браузер).
func ResolveReference(base *url.URL, ref string) (*url.URL, error) {
	if ref == "" {
		return nil, fmt.Errorf("empty reference")
	}

	// Парсим ссылку как относительный или абсолютный URL
	refURL, err := url.Parse(ref)
	if err != nil {
		return nil, fmt.Errorf("invalid reference URL %q: %w", ref, err)
	}

	absolute := base.ResolveReference(refURL)

	// Проверим, что результат имеет схему и хост
	if absolute.Scheme == "" || absolute.Host == "" {
		return nil, fmt.Errorf("resolved URL is not absolute: %s", absolute.String())
	}

	return absolute, nil
}
