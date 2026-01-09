package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Result содержит результат загрузки
type Result struct {
	URL  *url.URL // финальный URL (после редиректов)
	Body []byte   // тело ответа
}

// Options задаёт параметры загрузчика
type Options struct {
	Timeout time.Duration
}

var DefaultOptions = Options{
	Timeout: 10 * time.Second,
}

// Downloader выполняет HTTP-загрузки
type Downloader struct {
	client *http.Client
}

// New создаёт новый загрузчик с заданными опциями
func New(opts Options) *Downloader {
	client := &http.Client{
		Timeout: opts.Timeout,
	}
	return &Downloader{client: client}
}

// Fetch загружает URL и возвращает тело ответа
func (d *Downloader) Fetch(u *url.URL) (*Result, error) {
	resp, err := d.client.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", u.String(), err)
	}
	defer resp.Body.Close()

	// Проверим статус
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got non-200 status: %d for %s", resp.StatusCode, u.String())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from %s: %w", u.String(), err)
	}

	// Финальный URL может отличаться из-за редиректов
	finalURL, err := url.Parse(resp.Request.URL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to parse url %s: %w", u.String(), err)
	}

	return &Result{
		URL:  finalURL,
		Body: body,
	}, nil
}
