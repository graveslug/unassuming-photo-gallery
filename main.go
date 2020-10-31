package main

import (
	"net/http"

	"github.com/graveslug/unassuming-photo-gallery/views"

	"github.com/gorilla/mux"
)

//Will remove these after. Globals suck.
var (
	homeView    *views.View
	contactView *views.View
	faqView     *views.View
	signupView  *views.View
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	errHand(homeView.Render(w, nil))

}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html")
	errHand(contactView.Render(w, nil))
}
func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	errHand(faqView.Render(w, nil))
}
func signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	errHand(signupView.Render(w, nil))
}

func main() {
	homeView = views.NewView("bootstrap", "views/home.gohtml")
	contactView = views.NewView("bootstrap", "views/contact.gohtml")
	faqView = views.NewView("bootstrap", "views/faq.gohtml")
	signupView = views.NewView("bootstrap", "views/signup.gohtml")

	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/signup", signup)
	r.HandleFunc("/faq", faq)
	r.HandleFunc("/contact", contact)
	http.ListenAndServe(":3000", r)
}

//errHand helper function that panics (like me) when there is an error
func errHand(err error) {
	if err != nil {
		panic(err)
	}
}
