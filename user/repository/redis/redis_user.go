package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis"
	"github.com/menduong/oauth2/domain"
)

type redisUserRepository struct {
	redis *redis.Client
}

// NewRedisUserRepository will create an object that represent the article.Repository interface
func NewRedisUserRepository(Conn *redis.Client) domain.UserOTPRepository {
	return &redisUserRepository{Conn}
}

// GeOTP is
func (m *redisUserRepository) GetOTP(ctx context.Context, email string) (string, error) {

	return m.redis.Get("otp:" + email).Result()
}

// SetOTP is
func (m *redisUserRepository) SetOTP(ctx context.Context, email string, otp string, expireTime time.Duration) error {
	m.redis.Set("otp:"+email, otp, expireTime)
	return nil
}
