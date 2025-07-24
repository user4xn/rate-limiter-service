package limiter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"rate-limiter/pkg/dto"
	"rate-limiter/pkg/factory"
	"rate-limiter/pkg/repository"
	"rate-limiter/pkg/util"
	"time"

	"github.com/go-redis/redis/v8"
)

type Service interface {
	FixedWindow(c context.Context, payload dto.PayloadLimiter) (*dto.ResponseLimiter, error)
	SetClientConfigFixedWindow(c context.Context, payload dto.PayloadConfigClient) error
}

type service struct {
	redisRepository repository.RedisRepositoryInterface
	requestLimit    int
	window          int
}

func NewService(f *factory.Factory) Service {
	return &service{
		redisRepository: f.RedisRepository,
		requestLimit:    util.LoadDefaultMaxRequest(),
		window:          util.LoadDefaultWindow(),
	}
}

func (s *service) FixedWindow(c context.Context, payload dto.PayloadLimiter) (*dto.ResponseLimiter, error) {
	var (
		res      dto.ResponseLimiter
		limit    int
		duration time.Duration
	)

	if payload.ClientID == "" {
		return nil, fmt.Errorf("client_id is required")
	}

	if payload.Route == "" {
		return nil, fmt.Errorf("route is required")
	}

	conf, err := s.GetClientConfig(c, payload.ClientID, payload.Route)
	if err != nil {
		return nil, fmt.Errorf("failed to get client config: %w", err)
	}

	if conf.Limit != 0 {
		limit = conf.Limit
	}
	if conf.Window != 0 {
		duration = time.Second * time.Duration(conf.Window)
	}

	now := time.Now().UTC()
	window := now.Truncate(duration).Unix()
	key := fmt.Sprintf("fixed-window:%s:%s:%d", payload.ClientID, payload.Route, window)

	// check increment if doest exist then create (1)
	current, err := s.redisRepository.Incr(c, key)
	if err != nil {
		return nil, fmt.Errorf("failed to increment counter: %w", err)
	}

	// set TTL only if current is 1
	if current == 1 {
		ttl := duration - now.Sub(time.Unix(window, 0))
		err = s.redisRepository.Expire(c, key, ttl)
		if err != nil {
			return nil, fmt.Errorf("failed to set TTL: %w", err)
		}
	}

	// if current is greater than request limit then return limit exceeded
	if current > int64(limit) {
		ttl, err := s.redisRepository.TTL(c, key)
		if err != nil {
			return nil, fmt.Errorf("failed to get TTL: %w", err)
		}
		res = dto.ResponseLimiter{
			Status:        http.StatusTooManyRequests,
			Limit:         limit,
			Remain:        0,
			ResetInSecond: int(ttl.Seconds()),
		}
		return &res, nil
	}

	// get TTL
	ttl, err := s.redisRepository.TTL(c, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get TTL: %w", err)
	}

	// return available remain and TTL
	res = dto.ResponseLimiter{
		Status:        http.StatusOK,
		Limit:         limit,
		Remain:        limit - int(current),
		ResetInSecond: int(ttl.Seconds()),
	}

	return &res, nil
}

// get client config adn save if not exist with default value
func (s *service) GetClientConfig(ctx context.Context, clientID string, route string) (dto.ConfigClientRedis, error) {
	key := fmt.Sprintf("fixed-window-config:%s:%s", clientID, route)

	val, err := s.redisRepository.Get(ctx, key)
	if err != nil && err != redis.Nil {
		return dto.ConfigClientRedis{}, err
	}

	if val == "" {
		payload := dto.ConfigClientRedis{
			ClientID: clientID,
			Limit:    s.requestLimit,
			Window:   s.window,
		}

		err = s.redisRepository.Set(ctx, key, payload, 24*time.Hour)
		if err != nil {
			return dto.ConfigClientRedis{}, err
		}

		return payload, nil
	}

	var conf dto.ConfigClientRedis
	err = json.Unmarshal([]byte(val), &conf)
	if err != nil {
		return dto.ConfigClientRedis{}, err
	}

	return conf, nil
}

func (s *service) SetClientConfigFixedWindow(c context.Context, payload dto.PayloadConfigClient) error {
	key := fmt.Sprintf("fixed-window-config:%s:%s", payload.ClientID, payload.Route)

	if payload.Limit == 0 {
		payload.Limit = s.requestLimit
	}

	if payload.Window == 0 {
		payload.Window = s.window
	}

	err := s.redisRepository.Set(c, key, payload, 24*time.Hour)
	if err != nil {
		return err
	}
	return nil
}
