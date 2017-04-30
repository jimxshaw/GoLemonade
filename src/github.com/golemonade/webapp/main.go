package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	// Get our list of templates.
	templates := populateTemplates()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Find the name of the template we're seeking. When using
		// a template parsing method like ParseGlob, the template
		// name is its file name.
		// Slice the first character off the path because it will
		// always be a slash /.
		requestedFile := r.URL.Path[1:]

		t := templates[requestedFile+".html"]
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

// We return a map of strings to templates. Previously we had a single template
// containing all other templates. Now we define a family of templates with each
// one individually cloned from _layout.html and then pulling in the content
// template that defines the look and feel for that template.
func populateTemplates() map[string]*template.Template {
	result := make(map[string]*template.Template)
	const basePath = "templates"
	layout := template.Must(template.ParseFiles(basePath + "/_layout.html"))
	template.Must(layout.ParseFiles(basePath+"/_header.html", basePath+"/_footer.html"))
	dir, err := os.Open(basePath + "/content")
	if err != nil {
		panic("Failed to open template blocks directory: " + err.Error())
	}
	fis, err := dir.Readdir(-1)
	if err != nil {
		panic("Failed to read contents of content directory: " + err.Error())
	}
	for _, fi := range fis {
		f, err := os.Open(basePath + "/content/" + fi.Name())
		if err != nil {
			panic("Failed to open template '" + fi.Name() + "'")
		}
		content, err := ioutil.ReadAll(f)
		if err != nil {
			panic("Failed to read content from file '" + fi.Name() + "'")
		}
		f.Close()
		tmpl := template.Must(layout.Clone())
		_, err = tmpl.Parse(string(content))
		if err != nil {
			panic("Failed to parse contents of '" + fi.Name() + "' as template")
		}
		result[fi.Name()] = tmpl
	}
	return result
}

// func populateTemplates() *template.Template {
// 	// Here's the container that has all of the templates we'll
// 	// be loading in.
// 	result := template.New("templates")

// 	// Location where the templates are stored on the file system.
// 	const basePath = "templates"

// 	// Parse the templates in the context of the result container.
// 	template.Must(result.ParseGlob(basePath + "/*.html"))

// 	// All the parsed out templates will be the children of this
// 	// result template container. That doesn't matter to Go as
// 	// it will look through the entire template tree when given
// 	// a template to execute.
// 	return result
// }
