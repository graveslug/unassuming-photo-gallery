package views

import (
	"html/template"
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

//NewView parses our view pages to reduce overall maintenance and repeat code in main.go
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
