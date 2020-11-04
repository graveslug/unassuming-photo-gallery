package controllers

import (
	"fmt"
	"net/http"

	"github.com/graveslug/unassuming-photo-gallery/models"
	"github.com/graveslug/unassuming-photo-gallery/views"
)

//NewUsers This will setup all the views we will need to handle for the user controller which will make it easier to reuse our controllers later on
func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
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
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "User is", user)
}

//Login is used to process the login form of a user when they try to login in as an existing user via email and PW
//POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user, err := u.us.Authenticate(form.Email, form.Password)
	switch err {
	case models.ErrNotFound:
		fmt.Fprintln(w, "Invalid email address.")
	case models.ErrInvalidPassword:
		fmt.Fprintln(w, "Invalid password Provided.")
	case nil:
		fmt.Fprintln(w, user)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//Users structure for the controller
type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        *models.UserService
}

//SignupForm repersents the input fields of our signup form.
type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

//LoginForm used as the parameters for logging in
type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}
