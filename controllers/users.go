package controllers

import (
	"fmt"
	"net/http"

	"github.com/graveslug/unassuming-photo-gallery/views"
)

//NewUsers This will setup all the views we will need to handle for the user controller which will make it easier to reuse our controllers later on
func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "users/new"),
	}
}

//New is used to render the form where the user can create a new user account.
//GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// Create route to process the signup form when the user tries to create a new user account
//POST /Signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, "Email is", form.Email)
	fmt.Fprintln(w, "Password is", form.Password)
}

//Users structure for the controller
type Users struct {
	NewView *views.View
}

//SignupForm repersents the input fields of our signup form.
type SignupForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}
