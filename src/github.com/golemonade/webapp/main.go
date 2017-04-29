package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Look into request object and see what file was requested and
		// then translate that into our file system.
		f, err := os.Open("public" + r.URL.Path)

		// Check to see if we successfully opened that file.
		if err != nil {
			// Write out a header.
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}

		// Assuming everything succeeded, we close it.
		defer f.Close()

		// Copy the file directly to the response writer.
		io.Copy(w, f)
	})

	http.ListenAndServe(":8000", nil)
}
