package routes

import (
	"BackendGoLdap/handlers"

	"github.com/go-chi/chi/v5"
)

func InitRoutes(router *chi.Mux) {
	router.Get("/", handlers.MainPage)
	router.Get("/authorized", handlers.AuthorizedPage)
}
