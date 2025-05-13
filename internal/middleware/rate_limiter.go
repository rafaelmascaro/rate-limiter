package middleware

import (
	"context"
	"net"
	"net/http"

	"github.com/rafaelmascaro/rate-limiter/internal/usecase"
)

type RateLimiterMiddleware struct {
	RateLimiterUseCase usecase.RateLimiterUseCase
}

func NewRateLimiterMiddleware(
	rateLimiterUseCase usecase.RateLimiterUseCase,
) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		RateLimiterUseCase: rateLimiterUseCase,
	}
}

func (m *RateLimiterMiddleware) Handler(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		token := r.Header.Get("API_KEY")

		allowed, err := m.RateLimiterUseCase.Allow(ctx, ip, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !allowed {
			http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
