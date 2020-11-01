package main

import (
	"net/http"

	"github.com/graveslug/unassuming-photo-gallery/controllers"
	"github.com/graveslug/unassuming-photo-gallery/views"

	"github.com/gorilla/mux"
)

//Will remove these after. Globals suck.
var (
	faqView *views.View
)

func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	errHand(faqView.Render(w, nil))
}

func main() {
	faqView = views.NewView("bootstrap", "views/faq.gohtml")
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers()

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.HandleFunc("/faq", faq).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	http.ListenAndServe(":3000", r)
}

//errHand helper function that panics (like me) when there is an error
func errHand(err error) {
	if err != nil {
		panic(err)
	}
}
