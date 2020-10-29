package views

import "html/template"

//NewView handles our view pages to reduce overall maintenance and repeat code
func NewView(files ...string) *View {
	files = append(files, "views/layouts/footer.gohtml")
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
	}
}

//View is to point at the template which will eventually point at the compiled template
type View struct {
	Template *template.Template
}
