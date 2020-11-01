package main

import (
	"net/http"

	"github.com/graveslug/unassuming-photo-gallery/controllers"

	"github.com/gorilla/mux"
)

func main() {
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers()

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.HandleFunc("/faq", staticC.Faq).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	http.ListenAndServe(":3000", r)
}

//errHand helper function that panics (like me) when there is an error
func errHand(err error) {
	if err != nil {
		panic(err)
	}
}
