package views

import "html/template"

//NewView handles our view pages to reduce overall maintenance and repeat code
func NewView(layout string, files ...string) *View {
	files = append(files,
		"views/layouts/footer.gohtml",
		"views/layouts/bootstrap.gohtml")
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
