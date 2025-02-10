//go:build !solution

package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var (
	urlMap = make(map[string]string)
	mutex  = &sync.Mutex{}
)

type GetURL struct {
	URL string `json:"url"`
}

type ResponseURL struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

func generateRandomCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func UrlShortenerHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusBadRequest)
		return
	}

	var url GetURL
	if err := json.NewDecoder(req.Body).Decode(&url); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if url.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for key, value := range urlMap {
		if value == url.URL {
			response := ResponseURL{
				URL: url.URL,
				Key: key,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	code := generateRandomCode()
	urlMap[code] = url.URL

	response := ResponseURL{
		URL: url.URL,
		Key: code,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func goHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/go/"):]

	mutex.Lock()
	url, ok := urlMap[key]
	mutex.Unlock()

	if !ok {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func main() {
	port := flag.String("port", "8081", "Port to run the server on")
	flag.Parse()

	http.HandleFunc("/shorten", UrlShortenerHandler)
	http.HandleFunc("/go/", goHandler)

	log.Printf("Server running on port %s...", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
