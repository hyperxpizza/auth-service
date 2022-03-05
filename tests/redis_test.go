package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/hyperxpizza/auth-service/pkg/auth"
	"github.com/hyperxpizza/auth-service/pkg/config"
	"github.com/stretchr/testify/assert"
)

//go test -v ./tests --run TestRedisConnection --config=/home/hyperxpizza/dev/golang/reusable-microservices/auth-service/config.json
func TestRedisConnection(t *testing.T) {
	flag.Parse()

	if *configPathOpt == "" {
		t.Fail()
		return
	}

	cfg, err := config.NewConfig(*configPathOpt)
	assert.NoError(t, err)

	rdc := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		DB:   int(cfg.Redis.DB),
	})

	stat, err := rdc.Ping(ctx).Result()
	assert.NoError(t, err)

	assert.Equal(t, stat, "PONG")
}

//go test -v ./tests --run TestRedisFunctions --config=/home/hyperxpizza/dev/golang/reusable-microservices/auth-service/config.json
func TestRedisFunctions(t *testing.T) {
	flag.Parse()

	if *configPathOpt == "" {
		t.Fail()
		return
	}

	cfg, err := config.NewConfig(*configPathOpt)
	assert.NoError(t, err)

	authenticator, err := auth.NewAuthenticator(cfg)
	assert.NoError(t, err)

	accessUID := "some-access-UID"
	refreshUID := "some-refresh-UID"

	cacheJSON, err := json.Marshal(auth.CachedTokens{
		AccessUID:  accessUID,
		RefreshUID: refreshUID,
	})
	assert.NoError(t, err)

	err = authenticator.SetTokensInRedis(*authServiceIDOpt, *usersServiceIDOpt, *usernameOpt, string(cacheJSON))
	assert.NoError(t, err)

	cj, err := authenticator.GetTokensFromRedis(*authServiceIDOpt, *usersServiceIDOpt, *usernameOpt)
	assert.NoError(t, err)

	ct := new(auth.CachedTokens)
	err = json.Unmarshal([]byte(cj), &ct)
	assert.NoError(t, err)

	assert.Equal(t, accessUID, ct.AccessUID)
	assert.Equal(t, refreshUID, ct.RefreshUID)

}
