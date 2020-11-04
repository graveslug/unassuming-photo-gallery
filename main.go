package main

import (
	"fmt"
	"net/http"

	"github.com/graveslug/unassuming-photo-gallery/controllers"
	"github.com/graveslug/unassuming-photo-gallery/models"

	"github.com/gorilla/mux"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "graveslug"
	dbname = "unassuming_photo_gallery"
)

func main() {
	//Create a DB connection string and then use it to create the model services
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.AutoMigrate()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(us)

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	//Why does this require the ServeHTTP while the others don't?!
	r.HandleFunc("/faq", staticC.Faq.ServeHTTP).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	http.ListenAndServe(":3000", r)
}

//errHand helper function that panics (like me) when there is an error
func errHand(err error) {
	if err != nil {
		panic(err)
	}
}
