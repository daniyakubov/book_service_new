package cache

import (
	"fmt"
	"github.com/daniyakubov/book_service_n/models"
	errors "github.com/fiverr/go_errors"
	"gopkg.in/redis.v5"
	"strings"
	"time"
)

var _ Cache = &RedisCache{}

type RedisCache struct {
	redisAdress string
	db          int
	expiration  time.Duration
	maxSize     int64
	client      *redis.Client
}

func NewRedisCache(redisAdress string, db int, exp time.Duration, maxSize int64, client *redis.Client) *RedisCache {
	return &RedisCache{
		redisAdress: redisAdress,
		db:          db,
		expiration:  exp,
		maxSize:     maxSize,
		client:      client,
	}
}
func (cache *RedisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.redisAdress,
		Password: "",
		DB:       cache.db,
	})
}

func (cache *RedisCache) AddAction(key string, method string, routeName string) error {
	length, err := cache.client.LLen(key).Result()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failded to get length of actions of username %s", key))
	}

	cache.client.RPush(key, fmt.Sprintf("method: %s, route: %s", method, routeName))
	if length >= cache.maxSize {
		_, err := cache.client.LPop(key).Result()
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failded to update action of username %s", key))
		}
	}

	return nil
}

func (cache *RedisCache) GetLastActions(key string) ([]models.Action, error) {
	length, err := cache.client.LLen(key).Result()
	if err != nil {
		panic(err)
	}
	items, err := cache.client.LRange(key, 0, length).Result()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failded to get actions of username %s", key))
	}

	res := make([]models.Action, len(items))
	for i := 0; i < len(items); i++ {
		s := strings.Split(items[i], ",")
		method := strings.Split(s[0], ":")[1]
		routeName := strings.Split(s[1], ":")[1]
		res[i].Method = method
		res[i].RouteName = routeName
	}

	return res, nil
}
