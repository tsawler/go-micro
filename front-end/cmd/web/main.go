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

//go:embed templates
var templateFS embed.FS

func render(w http.ResponseWriter, t string) {

	templates := []string{
		fmt.Sprintf("templates/%s", t), // order matters; page must come first
		"templates/header.partial.gohtml",
		"templates/footer.partial.gohtml",
		"templates/base.layout.gohtml",
	}

	tmpl, err := template.ParseFS(templateFS, templates...)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data struct {
		BrokerURL string
	}

	data.BrokerURL = os.Getenv("BROKER_URL")

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
