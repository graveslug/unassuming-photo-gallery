package controllers

import "github.com/graveslug/unassuming-photo-gallery/views"

//NewStatic is a handler that handles "static" pages rather than convulting the file base with one off controllers.
func NewStatic() *Static {
	return &Static{
		Home: views.NewView(
			"bootstrap", "views/static/home.gohtml"),
		Contact: views.NewView(
			"bootstrap", "views/static/contact.gohtml"),
	}
}

//Static structure to
type Static struct {
	Home    *views.View
	Contact *views.View
}
