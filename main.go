package main

import (
	"net/http"

	"github.com/graveslug/unassuming-photo-gallery/views"

	"github.com/gorilla/mux"
)

//Will remove these after. Globals suck.
var homeView *views.View
var contactView *views.View

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	errHand(homeView.Render(w, nil))

}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html")
	errHand(contactView.Render(w, nil))
}

func main() {
	//There is a an error with the bootstrap partt that causes an undefined views/contact.gohtml is undefined
	homeView = views.NewView("bootstrap", "views/home.gohtml")
	contactView = views.NewView("bootstrap", "views/contact.gohtml")

	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	http.ListenAndServe(":3000", r)
}

//errHand helper function that panics (like me) when there is an error
func errHand(err error) {
	if err != nil {
		panic(err)
	}
}
