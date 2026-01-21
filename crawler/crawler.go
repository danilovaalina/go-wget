package main

import (
	"fmt"
	"net/url"
	"sync"
)

type Crawler struct {
	downloader *Downloader
	saver      *Saver
	startHost  string
	visited    map[string]bool
	mu         sync.Mutex
}

func New(dl *Downloader, sv *Saver, startURL *url.URL) *Crawler {
	return &Crawler{
		downloader: dl,
		saver:      sv,
		startHost:  startURL.Host,
		visited:    make(map[string]bool),
	}
}

func (c *Crawler) isVisited(u *url.URL) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.visited[u.String()]
}

func (c *Crawler) markVisited(u *url.URL) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.visited[u.String()] = true
}

func (c *Crawler) isSameDomain(u *url.URL) bool {
	return u.Host == c.startHost
}

type Job struct {
	URL   *url.URL
	Depth int
}

func (c *Crawler) Crawl(startURL *url.URL, maxDepth int) error {
	if maxDepth < 0 {
		return nil
	}

	queue := []Job{{URL: startURL, Depth: maxDepth}}

	for len(queue) > 0 {
		job := queue[0]
		queue = queue[1:]

		if job.Depth < 0 {
			continue
		}

		if c.isVisited(job.URL) {
			continue
		}

		if !c.isSameDomain(job.URL) {
			continue
		}

		c.markVisited(job.URL)

		// Скачиваем
		result, err := c.downloader.Fetch(job.URL)
		if err != nil {
			fmt.Printf("⚠️ Skip %s: %v\n", job.URL, err)
			continue
		}

		// Сохраняем
		_, err = c.saver.Save(result.URL, result.Body)
		if err != nil {
			fmt.Printf("⚠️ Failed to save %s: %v\n", result.URL, err)
			continue
		}

		fmt.Printf("✅ Saved: %s (depth=%d)\n", result.URL, job.Depth)

		// Если это HTML и глубина > 0, парсим и добавляем новые ссылки
		if job.Depth > 0 && isHTML(result.URL) {
			resources, err := ParseHTML(result.Body)
			if err != nil {
				fmt.Printf("⚠️ Parse failed for %s: %v\n", result.URL, err)
				continue
			}

			allRefs := GetAllResources()
			for _, ref := range allRefs {
				absURL, err := urlutil.Resolve(result.URL, ref)
				if err != nil {
					continue // пропускаем невалидные ссылки
				}

				// Добавляем в очередь с уменьшенной глубиной
				queue = append(queue, Job{
					URL:   absURL,
					Depth: job.Depth - 1,
				})
			}
		}
	}

	return nil
}

// isHTML — простая эвристика: проверяем расширение или Content-Type позже можно улучшить
func isHTML(u *url.URL) bool {
	path := u.Path
	return path == "" ||
		path == "/" ||
		path[len(path)-1] == '/' ||
		urlutil.HasHTMLExtension(path)
}
