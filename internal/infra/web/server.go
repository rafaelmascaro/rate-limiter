package web

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rafaelmascaro/Rate-Limiter/internal/middleware"
)

type Webserver struct {
	RateLimiterMiddleware middleware.RateLimiterMiddleware
}

func NewServer(
	rateLimiterMiddleware middleware.RateLimiterMiddleware,
) *Webserver {
	return &Webserver{
		RateLimiterMiddleware: rateLimiterMiddleware,
	}
}

func (we *Webserver) CreateServer() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", we.RateLimiterMiddleware.Handler(we.HandleRequest))
	return router
}

func (we *Webserver) HandleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}
