package register

import (
	"log"
	"net/http"
	"strings"

	"github.com/bradacina/inventory/deps"
	"github.com/bradacina/inventory/httphelp"
)

const (
	fieldPassword  = "password"
	fieldPassword2 = "password2"
	fieldEmail     = "email"
	fieldEmail2    = "email2"

	errorMissingField       = "Missing required field:"
	errorFieldsDoNotMatch   = "Fields do not match:"
	errorPasswordLength     = "Password must have at least 8 characters"
	errorPasswordComplexity = "Password must contain at least one letter, one number and one symbol"

	alphanumeric = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"
	symbols      = "`~!@#$%^&*()_+-=<>,./?:;'\""
)

type RegisterHandler struct {
	*deps.Deps
}

type registerForm struct {
	Password  string
	Password2 string
	Email     string
	Email2    string
}

type registerFormFeedback struct {
	registerForm
	Errors []string
}

func NewHandler(deps *deps.Deps) *RegisterHandler {
	return &RegisterHandler{deps}
}

func (rh *RegisterHandler) Register(w http.ResponseWriter, r *http.Request) {
	rh.CookieAuth.DeleteLoginCookie(w)
	if r.Method == http.MethodGet {
		httphelp.ServeTemplate(w, httphelp.TemplateRegister, nil)

	} else if r.Method == http.MethodPost {
		regForm := registrationForm(r)
		if valid := registrationValidation(regForm); valid != nil {
			httphelp.ServeTemplate(w, httphelp.TemplateRegister, valid)
			return
		}

		err := rh.UserService.RegisterUser(regForm.Email, regForm.Password)
		if err != nil {
			feedback := registerFormFeedback{Errors: []string{"That email is already in use"}}
			httphelp.ServeTemplate(w, httphelp.TemplateRegister, feedback)
		}

		log.Println("New Registration:", regForm)
		// TODO: redirect to account
	} else {
		http.NotFound(w, r)
	}
}

func registrationForm(r *http.Request) *registerForm {
	var regForm registerForm
	regForm.Email = r.FormValue(fieldEmail)
	regForm.Email2 = r.FormValue(fieldEmail2)
	regForm.Password = r.FormValue(fieldPassword)
	regForm.Password2 = r.FormValue(fieldPassword2)
	return &regForm
}

func registrationValidation(regForm *registerForm) *registerFormFeedback {

	// TODO: don't copy over values if the form is valid
	regFeedback := registerFormFeedback{}
	regFeedback.Email = regForm.Email
	regFeedback.Email2 = regForm.Email2

	if len(regForm.Email) == 0 {
		regFeedback.Errors = append(regFeedback.Errors, errorMissingField+"Email")
	}
	if len(regForm.Email2) == 0 {
		regFeedback.Errors = append(regFeedback.Errors, errorMissingField+"Confirm Email")
	}
	if len(regForm.Password) == 0 {
		regFeedback.Errors = append(regFeedback.Errors, errorMissingField+"Password")
	}
	if len(regForm.Password2) == 0 {
		regFeedback.Errors = append(regFeedback.Errors, errorMissingField+"Confirm Password")
	}

	if len(regFeedback.Errors) > 0 {
		return &regFeedback
	}

	// fields match the confirmation fields

	if regForm.Email != regForm.Email2 {
		regFeedback.Errors = append(regFeedback.Errors, errorFieldsDoNotMatch+
			"Email and Confirm Email")
	}
	if regForm.Password != regForm.Password2 {
		regFeedback.Errors = append(regFeedback.Errors, errorFieldsDoNotMatch+
			"Password and Confirm Password")
	}

	if len(regFeedback.Errors) > 0 {
		return &regFeedback
	}

	// password complexity

	if len(regForm.Password) < 8 {
		regFeedback.Errors = append(regFeedback.Errors, errorPasswordLength)
	}

	if !strings.ContainsAny(regForm.Password, alphanumeric) ||
		!strings.ContainsAny(regForm.Password, symbols) {
		regFeedback.Errors = append(regFeedback.Errors, errorPasswordComplexity)
	}

	if len(regFeedback.Errors) > 0 {
		return &regFeedback
	}

	return nil
}
