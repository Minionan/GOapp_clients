// main.go
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"BBCapp/code/database"
	"BBCapp/code/pages"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("pages/*.html"))

	if err := database.Initialize("./db/data.db"); err != nil {
		log.Fatal(err)
	}
}

func main() {
	defer database.Close()

	// Routes
	http.HandleFunc("/", pages.ClientsHandler(tpl))
	http.HandleFunc("/clients/add", pages.ClientAddHandler(tpl))
	http.HandleFunc("/clients/edit/", pages.ClientEditHandler(tpl))
	http.HandleFunc("/clients/delete/", pages.ClientDeleteHandler())
	http.HandleFunc("/clients/export", pages.ClientExportHandler(tpl))
	http.HandleFunc("/clients/import", pages.ClientImportHandler())

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
