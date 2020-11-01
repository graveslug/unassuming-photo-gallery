package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

//this is setting up the pattern for the Glob pattern
var (
	LayoutDir   string = "views/layouts/"
	TemplateDir string = "views/"
	TemplateExt string = ".gohtml"
)

//addTemplatePath takes a slice of strings representing file paths for the template, and it prepends the TemplateDir to each string in the slice.
//eg the input of "home" would look like "views/home" if TemplateDir == "views/"
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

//addTemplateExt takes a slice of strings repersenting file paths for the template and appends the TemplateExt to each string of the slice
//eg the input of "home" would result in the output of "home.gohtml" if templateExt == ".gohtml"
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}

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
	addTemplatePath(files)
	addTemplateExt(files)
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
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := v.Render(w, nil); err != nil {
		panic(err)
	}
}
