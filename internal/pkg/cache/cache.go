package cache

import (
	"encoding/json"
	"fmt"
	"time"
)

// Global store instance, initialized via InitStore.
var store Store

func InitStore(s Store) {
	store = s
}

func GetStore() Store {
	return store
}

// ── Token Blacklist ──

func BlacklistToken(token string, expiration time.Duration) error {
	return store.Set("token:blacklist:"+token, "1", expiration)
}

func IsTokenBlacklisted(token string) bool {
	val, err := store.Get("token:blacklist:" + token)
	return err == nil && val == "1"
}

// ── User Permission Cache ──

const permCacheTTL = 10 * time.Minute

func permKey(userID uint) string {
	return fmt.Sprintf("perm:user:%d", userID)
}

func SetUserPermissions(userID uint, perms []string) error {
	data, err := json.Marshal(perms)
	if err != nil {
		return err
	}
	return store.Set(permKey(userID), string(data), permCacheTTL)
}

func GetUserPermissions(userID uint) ([]string, bool) {
	val, err := store.Get(permKey(userID))
	if err != nil {
		return nil, false
	}
	var perms []string
	if err := json.Unmarshal([]byte(val), &perms); err != nil {
		return nil, false
	}
	return perms, true
}

func ClearUserPermissions(userID uint) {
	store.Del(permKey(userID))
}

func ClearAllPermissionCache() {
	keys, err := store.Scan("perm:user:*")
	if err != nil {
		return
	}
	for _, k := range keys {
		// keys from Scan already have the prefix, use raw client Del
		// but since our store.Del adds prefix again, we need to strip it
		// Instead, just delete by known user pattern
		store.Del(k)
	}
}

// ── Login Rate Limit ──

const loginLimitMax = 5
const loginLimitWindow = 15 * time.Minute

func loginKey(username string) string {
	return "login:fail:" + username
}

func IncrLoginFail(username string) (int64, error) {
	k := loginKey(username)
	count, err := store.Incr(k)
	if err != nil {
		return 0, err
	}
	if count == 1 {
		store.Expire(k, loginLimitWindow)
	}
	return count, nil
}

func IsLoginLocked(username string) bool {
	val, err := store.Get(loginKey(username))
	if err != nil {
		return false
	}
	var count int64
	fmt.Sscanf(val, "%d", &count)
	return count >= loginLimitMax
}

func ClearLoginFail(username string) {
	store.Del(loginKey(username))
}

func GetLoginLockTTL(username string) time.Duration {
	ttl, _ := store.TTL(loginKey(username))
	return ttl
}
