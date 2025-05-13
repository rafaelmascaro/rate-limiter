package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rafaelmascaro/rate-limiter/configs"
	"github.com/rafaelmascaro/rate-limiter/internal/entity"
	"github.com/rafaelmascaro/rate-limiter/internal/repository"
	"github.com/rafaelmascaro/rate-limiter/internal/usecase"
	"github.com/stretchr/testify/assert"
)

const (
	dbTest               int = 4
	rateLimitDefaultTest int = 10
	timeBlockDefaultTest int = 20
)

func SetupServerTest() (entity.RateLimiterRepository, *httptest.Server) {
	configs, err := configs.LoadConfig("../..")
	if err != nil {
		panic(err)
	}

	rateLimiterRepository := repository.NewRedisRepository(
		configs.RedisHost,
		configs.RedisPort,
		dbTest,
	)
	rateLimiterUseCase := usecase.NewRateLimiterUseCase(
		rateLimiterRepository,
		rateLimitDefaultTest,
		timeBlockDefaultTest,
	)
	rateLimiterMiddleware := NewRateLimiterMiddleware(*rateLimiterUseCase)

	mux := http.NewServeMux()
	mux.Handle("/", rateLimiterMiddleware.Handler(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}))

	return rateLimiterRepository, httptest.NewServer(mux)
}

func AddRateLimitForTest(
	repo entity.RateLimiterRepository,
	identifier string,
	limit int,
	timeBlock time.Duration,
) {
	key := fmt.Sprintf("rate_limit:%s", identifier)
	ctx := context.Background()
	repo.AddHash(ctx, key, limit, timeBlock)
}

func DeleteRateLimitForTest(repo entity.RateLimiterRepository, identifier string) {
	key := fmt.Sprintf("rate_limit:%s", identifier)
	ctx := context.Background()
	repo.Delete(ctx, key)
}

func TestRateLimiterDefault(t *testing.T) {
	_, server := SetupServerTest()
	defer server.Close()

	limit := 10
	secondsBlock := 20
	secondsVerify := 5

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/", nil)

	for i := 0; i < limit; i++ {
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}

	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)

	time.Sleep(time.Duration(secondsVerify) * time.Second)

	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)

	time.Sleep((time.Duration(secondsBlock - secondsVerify)) * time.Second)

	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRateLimiterToken(t *testing.T) {
	repo, server := SetupServerTest()
	defer server.Close()

	limit := 10
	secondsBlock := 10
	secondsVerify := 5

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/", nil)
	req.Header.Set("API_KEY", "abc123")

	AddRateLimitForTest(repo, "127.0.0.1", 5, 15*time.Second)

	for i := 0; i < limit; i++ {
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}

	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)

	time.Sleep(time.Duration(secondsVerify) * time.Second)

	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)

	time.Sleep((time.Duration(secondsBlock - secondsVerify)) * time.Second)

	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	DeleteRateLimitForTest(repo, "127.0.0.1")
}

func TestRateLimiterIP(t *testing.T) {
	repo, server := SetupServerTest()
	defer server.Close()

	limit := 5
	secondsBlock := 15
	secondsVerify := 5

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/", nil)

	AddRateLimitForTest(repo, "127.0.0.1", 5, 15*time.Second)

	for i := 0; i < limit; i++ {
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}

	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)

	time.Sleep(time.Duration(secondsVerify) * time.Second)

	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)

	time.Sleep((time.Duration(secondsBlock - secondsVerify)) * time.Second)

	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	DeleteRateLimitForTest(repo, "127.0.0.1")
}
