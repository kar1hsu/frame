package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisStore implements Store using go-redis.
// All keys are automatically prefixed with the configured prefix.
type RedisStore struct {
	client *redis.Client
	prefix string
}

func NewRedisStore(client *redis.Client, prefix string) *RedisStore {
	return &RedisStore{client: client, prefix: prefix}
}

func (s *RedisStore) ctx() context.Context {
	return context.Background()
}

func (s *RedisStore) key(k string) string {
	return s.prefix + k
}

// ── String ──

func (s *RedisStore) Get(key string) (string, error) {
	return s.client.Get(s.ctx(), s.key(key)).Result()
}

func (s *RedisStore) Set(key string, value string, expiration time.Duration) error {
	return s.client.Set(s.ctx(), s.key(key), value, expiration).Err()
}

func (s *RedisStore) Del(key string) error {
	return s.client.Del(s.ctx(), s.key(key)).Err()
}

func (s *RedisStore) Exists(key string) (bool, error) {
	n, err := s.client.Exists(s.ctx(), s.key(key)).Result()
	return n > 0, err
}

func (s *RedisStore) Incr(key string) (int64, error) {
	return s.client.Incr(s.ctx(), s.key(key)).Result()
}

func (s *RedisStore) Decr(key string) (int64, error) {
	return s.client.Decr(s.ctx(), s.key(key)).Result()
}

func (s *RedisStore) Expire(key string, expiration time.Duration) error {
	return s.client.Expire(s.ctx(), s.key(key), expiration).Err()
}

func (s *RedisStore) TTL(key string) (time.Duration, error) {
	return s.client.TTL(s.ctx(), s.key(key)).Result()
}

func (s *RedisStore) Scan(pattern string) ([]string, error) {
	ctx := s.ctx()
	var keys []string
	prefixLen := len(s.prefix)
	iter := s.client.Scan(ctx, 0, s.key(pattern), 100).Iterator()
	for iter.Next(ctx) {
		raw := iter.Val()
		if len(raw) > prefixLen {
			raw = raw[prefixLen:]
		}
		keys = append(keys, raw)
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}
	return keys, nil
}

// ── Hash ──

func (s *RedisStore) HGet(key, field string) (string, error) {
	return s.client.HGet(s.ctx(), s.key(key), field).Result()
}

func (s *RedisStore) HSet(key string, values ...interface{}) error {
	return s.client.HSet(s.ctx(), s.key(key), values...).Err()
}

func (s *RedisStore) HDel(key string, fields ...string) error {
	return s.client.HDel(s.ctx(), s.key(key), fields...).Err()
}

func (s *RedisStore) HGetAll(key string) (map[string]string, error) {
	return s.client.HGetAll(s.ctx(), s.key(key)).Result()
}

func (s *RedisStore) HExists(key, field string) (bool, error) {
	return s.client.HExists(s.ctx(), s.key(key), field).Result()
}

func (s *RedisStore) HIncrBy(key, field string, incr int64) (int64, error) {
	return s.client.HIncrBy(s.ctx(), s.key(key), field, incr).Result()
}

func (s *RedisStore) HKeys(key string) ([]string, error) {
	return s.client.HKeys(s.ctx(), s.key(key)).Result()
}

func (s *RedisStore) HLen(key string) (int64, error) {
	return s.client.HLen(s.ctx(), s.key(key)).Result()
}

func (s *RedisStore) HMGet(key string, fields ...string) ([]interface{}, error) {
	return s.client.HMGet(s.ctx(), s.key(key), fields...).Result()
}

// ── List ──

func (s *RedisStore) LPush(key string, values ...interface{}) (int64, error) {
	return s.client.LPush(s.ctx(), s.key(key), values...).Result()
}

func (s *RedisStore) RPush(key string, values ...interface{}) (int64, error) {
	return s.client.RPush(s.ctx(), s.key(key), values...).Result()
}

func (s *RedisStore) LPop(key string) (string, error) {
	return s.client.LPop(s.ctx(), s.key(key)).Result()
}

func (s *RedisStore) RPop(key string) (string, error) {
	return s.client.RPop(s.ctx(), s.key(key)).Result()
}

func (s *RedisStore) LRange(key string, start, stop int64) ([]string, error) {
	return s.client.LRange(s.ctx(), s.key(key), start, stop).Result()
}

func (s *RedisStore) LLen(key string) (int64, error) {
	return s.client.LLen(s.ctx(), s.key(key)).Result()
}

func (s *RedisStore) LRem(key string, count int64, value interface{}) (int64, error) {
	return s.client.LRem(s.ctx(), s.key(key), count, value).Result()
}

func (s *RedisStore) LIndex(key string, index int64) (string, error) {
	return s.client.LIndex(s.ctx(), s.key(key), index).Result()
}

func (s *RedisStore) LTrim(key string, start, stop int64) error {
	return s.client.LTrim(s.ctx(), s.key(key), start, stop).Err()
}

// ── Set ──

func (s *RedisStore) SAdd(key string, members ...interface{}) (int64, error) {
	return s.client.SAdd(s.ctx(), s.key(key), members...).Result()
}

func (s *RedisStore) SRem(key string, members ...interface{}) (int64, error) {
	return s.client.SRem(s.ctx(), s.key(key), members...).Result()
}

func (s *RedisStore) SMembers(key string) ([]string, error) {
	return s.client.SMembers(s.ctx(), s.key(key)).Result()
}

func (s *RedisStore) SIsMember(key string, member interface{}) (bool, error) {
	return s.client.SIsMember(s.ctx(), s.key(key), member).Result()
}

func (s *RedisStore) SCard(key string) (int64, error) {
	return s.client.SCard(s.ctx(), s.key(key)).Result()
}
