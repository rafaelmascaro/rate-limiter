package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rafaelmascaro/rate-limiter/configs"
	"github.com/rafaelmascaro/rate-limiter/internal/infra/web"
	"github.com/rafaelmascaro/rate-limiter/internal/middleware"
	"github.com/rafaelmascaro/rate-limiter/internal/repository"
	"github.com/rafaelmascaro/rate-limiter/internal/usecase"
)

func main() {
	configs, err := configs.LoadConfig("../..")
	if err != nil {
		panic(err)
	}

	rateLimiterRepository := repository.NewRedisRepository(
		configs.RedisHost,
		configs.RedisPort,
		configs.RedisDb,
	)
	rateLimiterUseCase := usecase.NewRateLimiterUseCase(
		rateLimiterRepository,
		configs.RateLimitDefault,
		configs.TimeBlockDefault,
	)
	rateLimiterMiddleware := middleware.NewRateLimiterMiddleware(
		*rateLimiterUseCase,
	)

	server := web.NewServer(*rateLimiterMiddleware)
	router := server.CreateServer()

	fmt.Println("Starting web server on port ", configs.WebServerPort)
	log.Fatal(http.ListenAndServe(configs.WebServerPort, router))
}
