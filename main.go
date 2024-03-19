package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
	"github.com/sajadjalilian/lenslocked/controllers"
	"github.com/sajadjalilian/lenslocked/migrations"
	"github.com/sajadjalilian/lenslocked/models"
	"github.com/sajadjalilian/lenslocked/templates"
	"github.com/sajadjalilian/lenslocked/views"
)

func main() {
	// Load ENV
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// host := os.Getenv("SMTP_HOST")
	// portStr := os.Getenv("SMTP_PORT")
	// port, err := strconv.Atoi(portStr)
	// if err != nil {
	// 	panic(err)
	// }
	// username := os.Getenv("SMTP_USERNAME")
	// password := os.Getenv("SMTP_PASSWORD")

	// Setup the database
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Setup Services
	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}

	// Setup Middleware
	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}
	csrfKey := "iyo7EQviU2kC*KXx8VpVjZgt9vnb9pxB"
	csrfMw := csrf.Protect([]byte(csrfKey), csrf.Secure(false))

	// Setup Controllers
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
	usersC.Templates.ForgetPassword =
		views.Must(views.ParseFS(
			templates.FS,
			"forgot-pw.gohtml",
			"tailwind.gohtml"))

	// Setup Router
	r := chi.NewRouter()
	r.Use(csrfMw)
	r.Use(umw.SetUser)
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
	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)
	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)
	r.Get("/users/me", usersC.CurrentUser)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	// Start the server
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
