package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

//this is setting up the pattern for the Glob pattern
var (
	LayoutDir   string = "views/layouts/"
	TemplateExt string = ".gohtml"
)

//function for the glob pattern. "views/layouts/*.gohtml"
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

//NewView parses our view pages to reduce overall maintenance and repeat code in main.go.
func NewView(layout string, files ...string) *View {
	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

//View is to point at the template which will eventually point at the compiled template
type View struct {
	Template *template.Template
	Layout   string
}

//Render this will the render method we will use instead of clogging up main.go with someView.Template.ExecuteTemplate(w, someView.Layout, nil)
func (v *View) Render(w http.ResponseWriter,
	data interface{}) error {
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}
