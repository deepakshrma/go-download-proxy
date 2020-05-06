package main

import (
	"io"
	"net/http"
	"os"
	"regexp"
	"sync"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
		// log.Fatal("$PORT must be set")
	}
	var validFiles = regexp.MustCompile(`\.(js|mp4|gz|zip|mp3|deb|dmg|csv|xml|tar|jar|apk|py|pptx|doc|docx)`)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var wg sync.WaitGroup

		url := r.URL.Query().Get("url")
		if validFiles.MatchString(url) {
			wg.Add(1)
			go func(url string) {
				println(url)
				// Decrement the counter when the go routine completes
				defer wg.Done()
				// Call the function check
				resp, err := http.Get(url)
				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()
				w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
				io.Copy(w, resp.Body)
			}(url)

		}
		// Wait for all the checkWebsite calls to finish
		wg.Wait()
	})
	println("Server is runing on:" + port)
	http.ListenAndServe(":"+port, nil)
}
