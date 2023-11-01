package main

import (
	"errors"
	"fmt"
	"net/http"

	"test.nicolas.net/pkg/models"
	"test.nicolas.net/pkg/validator"
)

type userSignupForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) homeDir(res http.ResponseWriter, req *http.Request) {
	rnd.HTML(res, http.StatusOK, "home", nil)
}

func (app *application) aboutDir(res http.ResponseWriter, req *http.Request) {
	rnd.HTML(res, http.StatusOK, "about", nil)
}

func (app *application) signInDir(res http.ResponseWriter, req *http.Request) {
	rnd.HTML(res, http.StatusOK, "signin", nil)
}

func (app *application) signUpDir(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "Display a HTML form for signing up a new user...")
}

func (app *application) signUpDirPost(res http.ResponseWriter, req *http.Request) {
	var form userSignupForm

	err := app.decodePostForm(req, &form)
	if err != nil {
		app.clientError(res, http.StatusBadRequest)
		return
	}

	form.CheckFields(validator.NotBlank(form.Email), "email", "Please fill this field!")
	form.CheckFields(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckFields(validator.NotBlank(form.Password), "password", "Please fill this field!")
	form.CheckFields(validator.MinChars(form.Password, 6), "password", "Password must be at least 6 characters long")

	if !form.Valid() {
		rnd.HTML(res, http.StatusOK, "signup", nil)
		return
	}

	err = app.users.Insert(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in used")

			rnd.HTML(res, http.StatusOK, "signup", nil)
		} else {
			app.serverError(res, err)
		}
		return

	}

	http.Redirect(res, req, "signin", http.StatusSeeOther)
}
