package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/sajadjalilian/lenslocked/controllers"
	"github.com/sajadjalilian/lenslocked/models"
	"github.com/sajadjalilian/lenslocked/templates"
	"github.com/sajadjalilian/lenslocked/views"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(
			templates.FS,
			"home.gohtml",
			"tailwind.gohtml"))))

	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(
			templates.FS,
			"contact.gohtml",
			"tailwind.gohtml"))))

	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS,
			"faq.gohtml",
			"tailwind.gohtml"))))

	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}

	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	usersC.Templates.New =
		views.Must(views.ParseFS(
			templates.FS,
			"signup.gohtml",
			"tailwind.gohtml"))

	usersC.Templates.SignIn =
		views.Must(views.ParseFS(
			templates.FS,
			"signin.gohtml",
			"tailwind.gohtml"))

	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Get("/users/me", usersC.CurrentUser)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	csrfKey := "iyo7EQviU2kC*KXx8VpVjZgt9vnb9pxB"
	csrfMw := csrf.Protect([]byte(csrfKey), csrf.Secure(false))

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", csrfMw(r))
}
