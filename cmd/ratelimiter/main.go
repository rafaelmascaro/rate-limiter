package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rafaelmascaro/Rate-Limiter/configs"
	"github.com/rafaelmascaro/Rate-Limiter/internal/infra/web"
	"github.com/rafaelmascaro/Rate-Limiter/internal/middleware"
	"github.com/rafaelmascaro/Rate-Limiter/internal/repository"
	"github.com/rafaelmascaro/Rate-Limiter/internal/usecase"
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
