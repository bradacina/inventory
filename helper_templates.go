package main

import "html/template"
import "log"
import "net/http"

type templateName string

const (
	TemplateLogin           templateName = "login.html"
	TemplateRegister        templateName = "register.html"
	TemplateInventoryList   templateName = "inventoryList.html"
	TemplateInventoryDetail templateName = "inventoryDetail.html"
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

func serveTemplate(w http.ResponseWriter, name templateName, values interface{}) {
	err := globalTemplate.ExecuteTemplate(w, string(name), values)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
