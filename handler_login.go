package main

import (
	"net/http"

	"github.com/bradacina/inventory/httphelp"
	"github.com/bradacina/inventory/logincookie"
)

type loginForm struct {
	Email    string
	Password string
}

type loginFormFeedback struct {
	loginForm
	Errors []string
}

func (app *app) login(w http.ResponseWriter, r *http.Request) {
	login, err := app.cookieAuth.GetLoginCookie(r)
	if err == nil {
		app.cookieAuth.SetLoginCookie(w, login)
		http.Redirect(w, r, "/secure", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		httphelp.ServeTemplate(w, httphelp.TemplateLogin, nil)
	} else if r.Method == http.MethodPost {
		form := parseLoginForm(r)

		if valid := validate(form); valid != nil {
			httphelp.ServeTemplate(w, httphelp.TemplateLogin, valid)
			return
		}

		user, err := app.userService.ValidateCredentials(form.Email, form.Password)
		if err != nil {
			feedback := loginFormFeedback{Errors: []string{err.Error()}}
			feedback.Email = form.Email
			httphelp.ServeTemplate(w, httphelp.TemplateLogin, feedback)
			return
		}

		if user.IsDeleted {
			feedback := loginFormFeedback{Errors: []string{"This user has been disabled"}}
			feedback.Email = form.Email
			httphelp.ServeTemplate(w, httphelp.TemplateLogin, feedback)
			return
		}

		loginInfo := logincookie.LoginInfo{Username: user.Email, ID: user.ID, IsAdmin: user.IsAdmin}
		app.cookieAuth.SetLoginCookie(w, &loginInfo)
		http.Redirect(w, r, "/secure", http.StatusSeeOther)

	} else {
		http.NotFound(w, r)
	}
}

func parseLoginForm(r *http.Request) *loginForm {
	var form loginForm
	form.Email = r.FormValue(fieldEmail)
	form.Password = r.FormValue(fieldPassword)
	return &form
}

func validate(form *loginForm) *loginFormFeedback {
	feedback := loginFormFeedback{}
	feedback.Email = form.Email

	if len(feedback.Email) == 0 || len(form.Password) == 0 {
		feedback.Errors = append(feedback.Errors, "You must enter an email and password")
	}

	if len(feedback.Errors) > 0 {
		return &feedback
	}

	return nil
}
