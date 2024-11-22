package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()

	var fr FileReader
	mux.HandleFunc("GET /posts/{slug}", PostHandler(fr))

	if err := http.ListenAndServe(":3030", mux); err != nil {
		log.Fatal(err)
	}
}

type SlugReader interface {
	Read(slug string) (string, error)
}

type FileReader struct{}

func (fr FileReader) Read(slug string) (string, error) {
	f, err := os.Open(slug + ".md")
	if err != nil {
		return "", err
	}

	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func PostHandler(sr SlugReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		postMarkdown, err := sr.Read(slug)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Write([]byte(postMarkdown))
	}
}
