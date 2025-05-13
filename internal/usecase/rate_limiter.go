package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/rafaelmascaro/Rate-Limiter/internal/entity"
)

type RateLimiterUseCase struct {
	Repo             entity.RateLimiterRepository
	RateLimitDefault int
	TimeBlockDefault time.Duration
}

func NewRateLimiterUseCase(
	repo entity.RateLimiterRepository,
	rateLimitDefault int,
	timeBlockDefault int,
) *RateLimiterUseCase {
	return &RateLimiterUseCase{
		Repo:             repo,
		RateLimitDefault: rateLimitDefault,
		TimeBlockDefault: time.Duration(timeBlockDefault) * time.Second,
	}
}

func (uc *RateLimiterUseCase) Allow(ctx context.Context, ip string, token string) (bool, error) {
	identifier, limit, timeBlock, err := uc.getRateLimit(ctx, ip, token)
	if err != nil {
		return false, err
	}

	blocked, err := uc.isBlocked(ctx, identifier)
	if err != nil {
		return false, err
	}
	if blocked {
		return false, nil
	}

	allowed, err := uc.isAllowedByCount(ctx, identifier, limit)
	if err != nil {
		return false, err
	}
	if !allowed {
		err = uc.block(ctx, identifier, timeBlock)
		return false, err
	}

	return true, nil
}

func (uc *RateLimiterUseCase) isAllowedByCount(ctx context.Context, identifier string, limit int) (bool, error) {
	key := fmt.Sprintf("count:%s", identifier)
	current, err := uc.Repo.Increment(ctx, key)
	if err != nil {
		return false, err
	}

	if current == 1 {
		err = uc.Repo.Expire(ctx, key, time.Second)
		if err != nil {
			return false, err
		}
	}

	return current <= int64(limit), nil
}

func (uc *RateLimiterUseCase) block(ctx context.Context, identifier string, timeBlock time.Duration) error {
	key := fmt.Sprintf("block:%s", identifier)
	return uc.Repo.AddKey(ctx, key, timeBlock)
}

func (uc *RateLimiterUseCase) isBlocked(ctx context.Context, identifier string) (bool, error) {
	key := fmt.Sprintf("block:%s", identifier)
	exists, err := uc.Repo.Exists(ctx, key)
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (uc *RateLimiterUseCase) getRateLimit(
	ctx context.Context,
	ip string,
	token string,
) (string, int, time.Duration, error) {
	identifier := token
	key := fmt.Sprintf("rate_limit:%s", identifier)
	exists, limit, timeBlock, err := uc.Repo.Find(ctx, key)
	if err != nil {
		return "", 0, 0, err
	}
	if exists {
		return identifier, limit, timeBlock, nil
	}

	identifier = ip
	key = fmt.Sprintf("rate_limit:%s", identifier)
	exists, limit, timeBlock, err = uc.Repo.Find(ctx, key)
	if err != nil {
		return "", 0, 0, err
	}
	if exists {
		return identifier, limit, timeBlock, nil
	}

	return identifier, uc.RateLimitDefault, uc.TimeBlockDefault, nil
}
