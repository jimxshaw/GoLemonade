package main

import (
	"html/template"
	"net/http"
)

func main() {
	http.ListenAndServe(":8000", http.FileServer(http.Dir("public")))
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
