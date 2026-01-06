package main

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kayden-vs/snippetbox/ui"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	// --- Global Middleware Stack ---
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(secureHeaders)

	// --- Public Routes ---
	staticFS, err := fs.Sub(ui.Files, "static")
	if err != nil {
		panic(err)
	}
	fileServer := http.FileServer(http.FS(staticFS))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	r.Get("/ping", ping)

	// --- Dynamic & Session-Enabled Routes ---
	r.Group(func(r chi.Router) {
		r.Use(app.sessionManager.LoadAndSave)
		r.Use(noSurf)
		r.Use(app.authenticate)

		r.Get("/", app.home)
		r.Get("/snippet/view/{id}", app.snippetView)
		r.Get("/user/signup", app.userSignup)
		r.Post("/user/signup", app.userSignupPost)
		r.Get("/user/login", app.userLogin)
		r.Post("/user/login", app.userLoginPost)
		r.Get("/about", app.about)

		// --- Authenticated-Only Routes ---
		r.Group(func(r chi.Router) {
			r.Use(app.requireAuthentication)

			r.Get("/snippet/create", app.snippetCreate)
			r.Post("/snippet/create", app.snippetCreatePost)
			r.Post("/user/logout", app.userLogoutPost)
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	return r
}
