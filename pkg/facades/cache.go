package facades

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
)

type redisCache struct {
	ctx      context.Context
	prefix   string
	instance *redis.Client
	store    string
}

func Cache() (*redisCache, error) {
	host := Env().GetString("REDIS_HOST", "127.0.0.1")
	if host == "" {
		return nil, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, Env().GetString("REDIS_PORT", "6379")),
		Password: Env().GetString("REDIS_PASSWORD"),
		DB:       Env().GetInt("REDIS_DATABASE", 0),
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, errors.WithMessage(err, "init connection error")
	}

	return &redisCache{
		ctx:      context.Background(),
		prefix:   fmt.Sprintf("%s:", Env().GetString("REDIS_PREFIX", "")),
		instance: client,
		store:    "",
	}, nil
}

// Add Driver an item in the cache if the key does not exist.
func (r *redisCache) Add(key string, value any, t time.Duration) bool {
	val, err := r.instance.SetNX(r.ctx, r.key(key), value, t).Result()
	if err != nil {
		return false
	}

	return val
}

func (r *redisCache) Decrement(key string, value ...int) (int, error) {
	if len(value) == 0 {
		value = append(value, 1)
	}

	res, err := r.instance.DecrBy(r.ctx, r.key(key), int64(value[0])).Result()

	return int(res), err
}

// Forever Driver an item in the cache indefinitely.
func (r *redisCache) Forever(key string, value any) bool {
	if err := r.Put(key, value, 0); err != nil {
		return false
	}

	return true
}

// Forget Remove an item from the cache.
func (r *redisCache) Forget(key string) bool {
	_, err := r.instance.Del(r.ctx, r.key(key)).Result()

	return err == nil
}

// Flush Remove all items from the cache.
func (r *redisCache) Flush() bool {
	res, err := r.instance.FlushAll(r.ctx).Result()

	if err != nil || res != "OK" {
		return false
	}

	return true
}

// Get Retrieve an item from the cache by key.
func (r *redisCache) Get(key string, def ...any) any {
	val, err := r.instance.Get(r.ctx, r.key(key)).Result()
	if err != nil {
		if len(def) == 0 {
			return nil
		}

		switch s := def[0].(type) {
		case func() any:
			return s()
		default:
			return s
		}
	}

	return val
}

func (r *redisCache) GetBool(key string, def ...bool) bool {
	if len(def) == 0 {
		def = append(def, false)
	}
	res := r.Get(key, def[0])
	if val, ok := res.(string); ok {
		return val == "1"
	}

	return cast.ToBool(res)
}

func (r *redisCache) GetInt(key string, def ...int) int {
	if len(def) == 0 {
		def = append(def, 1)
	}
	res := r.Get(key, def[0])
	if val, ok := res.(string); ok {
		i, err := strconv.Atoi(val)
		if err != nil {
			return def[0]
		}

		return i
	}

	return cast.ToInt(res)
}

func (r *redisCache) GetInt64(key string, def ...int64) int64 {
	if len(def) == 0 {
		def = append(def, 1)
	}
	res := r.Get(key, def[0])
	if val, ok := res.(string); ok {
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return def[0]
		}

		return i
	}

	return cast.ToInt64(res)
}

func (r *redisCache) GetString(key string, def ...string) string {
	if len(def) == 0 {
		def = append(def, "")
	}
	return cast.ToString(r.Get(key, def[0]))
}

// Has Check an item exists in the cache.
func (r *redisCache) Has(key string) bool {
	value, err := r.instance.Exists(r.ctx, r.key(key)).Result()

	if err != nil || value == 0 {
		return false
	}

	return true
}

func (r *redisCache) Increment(key string, value ...int) (int, error) {
	if len(value) == 0 {
		value = append(value, 1)
	}

	res, err := r.instance.IncrBy(r.ctx, r.key(key), int64(value[0])).Result()

	return int(res), err
}

// Put Driver an item in the cache for a given time.
func (r *redisCache) Put(key string, value any, t time.Duration) error {
	err := r.instance.Set(r.ctx, r.key(key), value, t).Err()
	if err != nil {
		return err
	}

	return nil
}

// Pull Retrieve an item from the cache and delete it.
func (r *redisCache) Pull(key string, def ...any) any {
	var res any
	if len(def) == 0 {
		res = r.Get(key)
	} else {
		res = r.Get(key, def[0])
	}
	r.Forget(key)

	return res
}

// Remember Get an item from the cache, or execute the given Closure and store the result.
func (r *redisCache) Remember(key string, seconds time.Duration, callback func() (any, error)) (any, error) {
	val := r.Get(key, nil)

	if val != nil {
		return val, nil
	}

	var err error
	val, err = callback()
	if err != nil {
		return nil, err
	}

	if err := r.Put(key, val, seconds); err != nil {
		return nil, err
	}

	return val, nil
}

// RememberForever Get an item from the cache, or execute the given Closure and store the result forever.
func (r *redisCache) RememberForever(key string, callback func() (any, error)) (any, error) {
	val := r.Get(key, nil)

	if val != nil {
		return val, nil
	}

	var err error
	val, err = callback()
	if err != nil {
		return nil, err
	}

	if err := r.Put(key, val, 0); err != nil {
		return nil, err
	}

	return val, nil
}

func (r *redisCache) key(key string) string {
	return r.prefix + key
}
