package main

import "html/template"
import "log"
import "net/http"

const (
	TemplateLogin    = "login.html"
	TemplateRegister = "register.html"
)

var (
	globalTemplate *template.Template
)

func init() {
	log.Println("Compiling templates")
	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	log.Println(tmpl.DefinedTemplates())

	globalTemplate = tmpl
}

func serveTemplate(w http.ResponseWriter, name string, values interface{}) {
	err := globalTemplate.ExecuteTemplate(w, name, values)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
