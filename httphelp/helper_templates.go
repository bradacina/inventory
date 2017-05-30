package httphelp

import "html/template"
import "log"
import "net/http"

type TemplateName string

const (
	TemplateLogin                TemplateName = "login.html"
	TemplateRegister             TemplateName = "register.html"
	TemplateInventoryList        TemplateName = "inventoryList.html"
	TemplateInventoryDetail      TemplateName = "inventoryDetail.html"
	TemplateAdminListUsers       TemplateName = "admin_list_users.html"
	TemplateAdminUser            TemplateName = "admin_user.html"
	TemplateAdminListInventories TemplateName = "admin_list_inventories.html"
	TemplateAdminInventory       TemplateName = "admin_inventory.html"
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

func ServeTemplate(w http.ResponseWriter, name TemplateName, values interface{}) {
	err := globalTemplate.ExecuteTemplate(w, string(name), values)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}
