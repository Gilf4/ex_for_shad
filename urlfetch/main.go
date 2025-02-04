//go:build !solution

package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func main() {
	urls := os.Args[1:]

	for _, v := range urls {
		parsedURL, err := url.ParseRequestURI(v)
		if err != nil {
			fmt.Printf("fetch: некорректный URL %s: %v\n", v, err)
			os.Exit(1)
		}

		resp, err := http.Get(parsedURL.String())
		if err != nil {
			fmt.Printf("fetch: ошибка при запросе к %s: %v\n", v, err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("fetch: ошибка при чтении ответа от %s: %v\n", v, err)
			os.Exit(1)
		}
	}
}
//all test passed
