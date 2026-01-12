package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/kayden-vs/snippetbox/internal/models"
	"github.com/kayden-vs/snippetbox/internal/validator"
	"github.com/kayden-vs/snippetbox/ui/html/pages"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.RenderPage(w, r, func(flash string, isAuthenticated bool, csrfToken string) templ.Component {
		return pages.HomePage(snippets, flash, isAuthenticated, csrfToken)
	})
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	authorName := snippet.AuthorName
	expiresStr := snippet.Expires.Format("02 Jan 2006 at 15:04")
	app.RenderPage(w, r, func(flash string, isAuthenticated bool, csrfToken string) templ.Component {
		return pages.ViewSnippet(snippet.ID, snippet.Title, snippet.Content, authorName, expiresStr, flash, isAuthenticated, csrfToken)
	})
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Create form instance with validator
	var form pages.SnippetCreateForm
	form.Validator = validator.Validator{FieldErrors: make(map[string]string)}

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Perform validation
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		app.RenderPage(w, r, func(flash string, isAuthenticated bool, csrfToken string) templ.Component {
			form.CSRFToken = csrfToken
			return pages.SnippetForm(form, isAuthenticated)
		})
		return
	}

	// Get the authenticated user's name
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	user, err := app.users.GetUserInfo(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	authorName := user.Name

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires, authorName)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.infoLog.Println("New Data added, Id: ", id)

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	form := pages.SnippetCreateForm{
		Expires:   365,
		Validator: validator.Validator{FieldErrors: make(map[string]string)},
	}

	app.RenderPage(w, r, func(flash string, isAuthenticated bool, csrfToken string) templ.Component {
		form.CSRFToken = csrfToken
		return pages.SnippetForm(form, isAuthenticated)
	})
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {

	props := pages.SignupFormParams{}
	app.RenderPage(w, r, func(flash string, isAuthenticated bool, csrfToken string) templ.Component {
		props.CSRFToken = csrfToken
		return pages.SignupPage(props, isAuthenticated)
	})
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form contents using our helper functions.
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	props := pages.SignupFormParams{
		Name:        form.Name,
		Email:       form.Email,
		FieldErrors: form.FieldErrors,
	}

	// If there are any errors, redisplay the signup form along with a 422 status code.
	if !form.Valid() {
		app.RenderPage(w, r, func(flash string, isAuthenticated bool, csrfToken string) templ.Component {
			props.CSRFToken = csrfToken
			return pages.SignupPage(props, isAuthenticated)
		})
		return
	}

	var id int

	id, err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			props.FieldErrors = form.FieldErrors
			app.RenderPage(w, r, func(flash string, isAuthenticated bool, csrfToken string) templ.Component {
				props.CSRFToken = csrfToken
				return pages.SignupPage(props, isAuthenticated)
			})
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Account created Succesfully.")

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	props := pages.LoginFormParams{}
	app.RenderPage(w, r, func(flash string, isAuthenticated bool, csrfToken string) templ.Component {
		props.CSRFToken = csrfToken
		return pages.LoginPage(props, flash, isAuthenticated)
	})
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	props := pages.LoginFormParams{
		Email:          form.Email,
		FieldErrors:    form.FieldErrors,
		NonFieldErrors: form.NonFieldErrors,
	}
	if !form.Valid() {
		app.RenderPage(w, r, func(flash string, isAuthenticated bool, csrfToken string) templ.Component {
			props.CSRFToken = csrfToken
			return pages.LoginPage(props, flash, isAuthenticated)
		})
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			props.NonFieldErrors = form.NonFieldErrors
			app.RenderPage(w, r, func(flash string, isAuthenticated bool, csrfToken string) templ.Component {
				props.CSRFToken = csrfToken
				return pages.LoginPage(props, flash, isAuthenticated)
			})
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out Succesfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	aboutContent := "CipherBin provides a clean, minimal interface for creating, viewing, and managing code snippets with automatic expiration. It features user authentication, session management, and a responsive design suitable for developers who need a quick way to share code samples."

	app.RenderPage(w, r, func(flash string, isAuthenticated bool, csrfToken string) templ.Component {
		return pages.AboutPage(aboutContent, csrfToken)
	})
}

func (app *application) viewAccount(w http.ResponseWriter, r *http.Request) {
	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	user, err := app.users.GetUserInfo(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.RenderPage(w, r, func(flash string, isAuthenticated bool, csrfToken string) templ.Component {
		return pages.AccountPage(user, flash, csrfToken)
	})
}

type passwordChangeForm struct {
	CurrentPassword     string `form:"currentPassword"`
	NewPassword         string `form:"newPassword"`
	ConfirmPassword     string `form:"newPasswordConfirmation"`
	validator.Validator `form:"-"`
}

func (app *application) updatePassword(w http.ResponseWriter, r *http.Request) {
	props := pages.PasswordFormParams{}
	app.RenderPage(w, r, func(flash string, isAuthenticated bool, csrfToken string) templ.Component {
		props.CSRFToken = csrfToken
		return pages.PasswordPage(props, isAuthenticated)
	})
}

func (app *application) updatePasswordPost(w http.ResponseWriter, r *http.Request) {
	var form passwordChangeForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//validation
	form.CheckField(validator.NotBlank(form.CurrentPassword), "currentPassword", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.NewPassword), "newPassword", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.ConfirmPassword), "confirmPassword", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.NewPassword, 8), "newPassword", "This field must be at least 8 characters long")
	form.CheckField(form.NewPassword == form.ConfirmPassword, "confirmPassword", "This should must be same as New password")

	props := pages.PasswordFormParams{
		CurrentPassword: form.CurrentPassword,
		NewPassword:     form.NewPassword,
		ConfirmPassword: form.ConfirmPassword,
		FieldErrors:     form.FieldErrors,
	}

	if !form.Valid() {
		app.RenderPage(w, r, func(flash string, isAuthenticated bool, csrfToken string) templ.Component {
			props.CSRFToken = csrfToken
			return pages.PasswordPage(props, isAuthenticated)
		})
		return
	}

	// compare password
	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	err = app.users.ComparePassword(id, form.CurrentPassword)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddFieldError("currentPassword", "Incorrect Password")
			props.FieldErrors = form.FieldErrors
			app.RenderPage(w, r, func(flash string, isAuthenticated bool, csrfToken string) templ.Component {
				props.CSRFToken = csrfToken
				return pages.PasswordPage(props, isAuthenticated)
			})
			return
		} else {
			app.serverError(w, err)
		}
	}

	// update password
	err = app.users.UpdatePassword(id, form.NewPassword)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Password updated successfully!")
	http.Redirect(w, r, "/account/view", http.StatusSeeOther)
}
