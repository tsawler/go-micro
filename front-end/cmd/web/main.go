package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.gohtml")
	})

	fmt.Println("Starting front end service on port 8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Panic(err)
	}
}

// tc is our template cache, a map of string/*template.Template.
// We look up our templates by the file name of the template.
var tc = make(map[string]*template.Template)

//go:embed templates
var templateFS embed.FS

func render(w http.ResponseWriter, t string) {
	// create two variables: one for the template to render, and an error
	var tmpl *template.Template
	var err error

	// check to see if we already have the template in the cache
	_, inMap := tc[t]
	if !inMap {
		// we don't have one, so create the template and add it to the cache
		err = createTemplateCache(t)
		if err != nil {
			// something went wrong
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	// pull the template out of the cache
	tmpl = tc[t]

	// create a struct so we can easily send data to the template we want to render
	var data struct {
		BrokerURL string
	}

	data.BrokerURL = os.Getenv("BROKER_URL")

	// execute the template, passing it data
	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// create template cache tries to parse the given template and add it to
// the template cache.
func createTemplateCache(t string) error {
	templates := []string{
		fmt.Sprintf("templates/%s", t),    // order matters; page must come first
		"templates/header.partial.gohtml", // these things, which the page depends on, can come in any order
		"templates/footer.partial.gohtml",
		"templates/base.layout.gohtml",
	}

	// parse the template
	tmpl, err := template.ParseFS(templateFS, templates...)
	if err != nil {
		return err
	}

	// add to map
	tc[t] = tmpl

	return nil
}
