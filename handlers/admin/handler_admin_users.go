package admin

import (
	"log"
	"net/http"
	"strconv"

	"github.com/bradacina/inventory/db"
	"github.com/bradacina/inventory/httphelp"
	"github.com/bradacina/inventory/routing"
	"github.com/bradacina/inventory/services"
)

type editUserFeedback struct {
	*db.User
	Errors []string
}

func (ah *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {

	users, err := ah.UserService.GetAll()
	if err != nil {
		httphelp.StatusCode(w, http.StatusInternalServerError)
		log.Panic(err)
	}
	httphelp.ServeTemplate(w, httphelp.TemplateAdminListUsers, users)
}

func (ah *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		httphelp.StatusCode(w, http.StatusInternalServerError)
	}

	requestUserID := r.FormValue("userId")
	userID, err := strconv.Atoi(requestUserID)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, routing.RouteAdminListUsers, http.StatusSeeOther)
		return
	}

	if userID == 1 {
		// cannot delete admin user
		http.Redirect(w, r, routing.RouteAdminListUsers, http.StatusSeeOther)
		return
	}

	user, err := ah.UserService.GetByID(userID)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, routing.RouteAdminListUsers, http.StatusSeeOther)
		return
	}

	user.IsDeleted = true
	err = ah.UserService.UpdateByAdmin(user)
	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, routing.RouteAdminListUsers, http.StatusSeeOther)
}

func (ah *AdminHandler) EditUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		r.ParseForm()

		userIDval := r.FormValue("userId")
		if len(userIDval) == 0 {
			httphelp.StatusCode(w, http.StatusNotFound)
			return
		}

		userID, err := strconv.Atoi(userIDval)
		if err != nil {
			log.Println(err)
			httphelp.StatusCode(w, http.StatusInternalServerError)
			return
		}

		user, err := ah.UserService.GetByID(userID)
		if err != nil {
			log.Println(err)
			httphelp.StatusCode(w, http.StatusNotFound)
			return
		}

		httphelp.ServeTemplate(w, httphelp.TemplateAdminUser, editUserFeedback{user, nil})
	} else if r.Method == http.MethodPost {
		user, err := parseUserFromRequest(r)
		if err != nil {
			httphelp.StatusCode(w, http.StatusInternalServerError)
			return
		}

		errors := ah.validateUserAndSave(user)
		if errors != nil && len(errors) > 0 {
			httphelp.ServeTemplate(w, httphelp.TemplateAdminUser, editUserFeedback{user, errors})
			return
		}

		http.Redirect(w, r, routing.RouteAdminListUsers, http.StatusSeeOther)
	}
}

func parseUserFromRequest(r *http.Request) (*db.User, error) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var user db.User

	err = decoder.Decode(&user, r.PostForm)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &user, nil
}

func (ah *AdminHandler) validateUserAndSave(user *db.User) []string {
	if user.ID == 1 && (!user.IsAdmin || user.IsDeleted) {
		return []string{"Cannot edit the super admin user"}
	}

	err := ah.UserService.UpdateByAdmin(user)
	if err == services.ErrorNotFound {
		return []string{"User not found in database"}
	}
	if err == services.ErrorEmailInUse {
		return []string{"That email address is already in use"}
	}
	if err != nil {
		log.Println(err)
		return nil
	}
	return nil
}
