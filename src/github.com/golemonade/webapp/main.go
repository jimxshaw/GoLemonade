package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	templates := populateTemplates()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Find the name of the template we're seeking. When using
		// a template parsing method like ParseGlob, the template
		// name is its file name.
		// Slice the first character off the path because it will
		// always be a slash /.
		requestedFile := r.URL.Path[1:]

		// Look up the template with the name of the requested file
		// and also adding the suffix. The user won't have to type
		// .html every time in the url if we include the suffix.
		t := templates.Lookup(requestedFile + ".html")
		if t != nil {
			err := t.Execute(w, nil)
			if err != nil {
				log.Println(err)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	// The way url resolution works in Go is that anything that comes
	// in prefixed with img or css will be handled by these two handlers
	// because they're more specific than the first handler above, which
	// only listens on the root url.
	http.Handle("/img/", http.FileServer(http.Dir("public")))
	http.Handle("/css/", http.FileServer(http.Dir("public")))

	http.ListenAndServe(":8000", nil)
}

func populateTemplates() *template.Template {
	// Here's the container that has all of the templates we'll
	// be loading in.
	result := template.New("templates")

	// Location where the templates are stored on the file system.
	const basePath = "templates"

	// Parse the templates in the context of the result container.
	template.Must(result.ParseGlob(basePath + "/*.html"))

	// All the parsed out templates will be the children of this
	// result template container. That doesn't matter to Go as
	// it will look through the entire template tree when given
	// a template to execute.
	return result
}
