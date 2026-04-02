package cache

import "time"

// Store defines a generic cache interface.
// Business code depends on this interface, not on go-redis directly.
type Store interface {
	// ── String ──
	Get(key string) (string, error)
	Set(key string, value string, expiration time.Duration) error
	Del(key string) error
	Exists(key string) (bool, error)
	Incr(key string) (int64, error)
	Decr(key string) (int64, error)
	Expire(key string, expiration time.Duration) error
	TTL(key string) (time.Duration, error)
	Scan(pattern string) ([]string, error)

	// ── Hash ──
	HGet(key, field string) (string, error)
	HSet(key string, values ...interface{}) error
	HDel(key string, fields ...string) error
	HGetAll(key string) (map[string]string, error)
	HExists(key, field string) (bool, error)
	HIncrBy(key, field string, incr int64) (int64, error)
	HKeys(key string) ([]string, error)
	HLen(key string) (int64, error)
	HMGet(key string, fields ...string) ([]interface{}, error)

	// ── List ──
	LPush(key string, values ...interface{}) (int64, error)
	RPush(key string, values ...interface{}) (int64, error)
	LPop(key string) (string, error)
	RPop(key string) (string, error)
	LRange(key string, start, stop int64) ([]string, error)
	LLen(key string) (int64, error)
	LRem(key string, count int64, value interface{}) (int64, error)
	LIndex(key string, index int64) (string, error)
	LTrim(key string, start, stop int64) error

	// ── Set ──
	SAdd(key string, members ...interface{}) (int64, error)
	SRem(key string, members ...interface{}) (int64, error)
	SMembers(key string) ([]string, error)
	SIsMember(key string, member interface{}) (bool, error)
	SCard(key string) (int64, error)
}
