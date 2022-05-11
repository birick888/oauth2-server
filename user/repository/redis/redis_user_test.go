package redis_test

import (
	"context"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis"
	_redisRepo "github.com/menduong/oauth2/user/repository/redis"
	"gopkg.in/go-playground/assert.v1"
)

var (
	redisServer *miniredis.Miniredis
	client      *redis.Client
)

// init data test
func init() {
}

func mockRedis() *miniredis.Miniredis {
	s, err := miniredis.Run()

	if err != nil {
		panic(err)
	}

	return s
}

func TestGetSetOTP(t *testing.T) {
	redisServer = mockRedis()
	client = redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})

	rd := _redisRepo.NewRedisUserRepository(client)
	err := rd.SetOTP(context.TODO(), "email1@gmail.com", "1234", time.Duration(5*time.Minute))
	if err != nil {
		t.Failed()
		return
	}
	otp, err := rd.GetOTP(context.TODO(), "email1@gmail.com")
	if err != nil {
		t.Failed()
		return
	}
	assert.Equal(t, otp, "1234")
}
