package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/graveslug/unassuming-photo-gallery/views"

	"github.com/gorilla/mux"
)

var homeTemplate *views.View
var contactTemplate *views.View

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := homeView.Template.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html")
	if err := contactTemplate.Execute(w, nil); err != nil {
		panic(err)
	}
}

func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>FAQ all the time for everytime</h1>")
}

func main() {
	var err error
	homeTemplate, err = template.ParseFiles(
		"views/home.gohtml")
	contactTemplate, err = template.ParseFiles(
		"views/contact.gohtml")

	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/faq", faq)
	r.HandleFunc("/contact", contact)
	http.ListenAndServe(":3000", r)
}
