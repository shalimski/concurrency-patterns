package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {
	urls := []string{
		"https://github.com",
		"https://go.dev",
		"https://badurl",
	}

	results := make([]string, len(urls))

	var wg sync.WaitGroup
	wg.Add(len(urls))
	for i, url := range urls {
		i := i
		go func(url string) {
			defer wg.Done()
			r, err := http.Get(url)
			if err != nil {
				results[i] = fmt.Sprintf("error: %v", err)
				return
			}
			results[i] = r.Status
		}(url)
	}
	wg.Wait()

	// ordered results
	for i, url := range urls {
		fmt.Printf("%s: %s\n", url, results[i])
	}
}
