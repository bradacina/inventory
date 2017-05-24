package main

import (
	"net/http"
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
	login, err := app.cookieHelper.getLoginCookie(r)
	if err == nil {
		app.cookieHelper.setLoginCookie(w, login)
		http.Redirect(w, r, "/secure", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		serveTemplate(w, TemplateLogin, nil)
	} else if r.Method == http.MethodPost {
		form := parseLoginForm(r)

		if valid := validate(form); valid != nil {
			serveTemplate(w, TemplateLogin, valid)
			return
		}

		user, err := app.userService.ValidateCredentials(form.Email, form.Password)
		if err != nil {
			feedback := loginFormFeedback{Errors: []string{err.Error()}}
			feedback.Email = form.Email
			serveTemplate(w, TemplateLogin, feedback)
			return
		}

		if user.IsDeleted {
			feedback := loginFormFeedback{Errors: []string{"This user has been disabled"}}
			feedback.Email = form.Email
			serveTemplate(w, TemplateLogin, feedback)
			return
		}

		loginInfo := loginInfo{Username: user.Email, ID: user.ID, IsAdmin: user.IsAdmin}
		app.cookieHelper.setLoginCookie(w, &loginInfo)
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
