package api

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenRepository struct {
	rdb *redis.Client
	ctx context.Context
}

func NewTokenRepository() ITokenRepo {
	opt, err := redis.ParseURL(config.getRedisUrl())
	if err != nil {
		log.Fatal(err)
	}
	client := redis.NewClient(opt)
	return &TokenRepository{rdb: client, ctx: nil}
}

func (t *TokenRepository) NewTokenRepositoryWithCtx(ctx context.Context) IRequestTokenRepo {
	return &TokenRepository{rdb: t.rdb, ctx: ctx}
}

func (t *TokenRepository) GetUserIdForToken(authToken string) (string, error) {
	val, err := t.rdb.Get(t.ctx, authToken).Result()
	return val, err
}

func (t *TokenRepository) RemoveToken(authToken string) error {
	return t.rdb.Del(t.ctx, authToken).Err()
}

func (t *TokenRepository) SetAuthToken(userId string, authToken string) error {
	return t.rdb.Set(t.ctx, authToken, userId, time.Minute*10).Err()
}
