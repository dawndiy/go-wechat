package token

import (
	"errors"
	"sync"
	"time"
)

var (
	// ErrTokenNotFound Store 中没有 Token
	ErrTokenNotFound = errors.New("token not found")
	// ErrTokenExpired Store 中的 Token 过期
	ErrTokenExpired = errors.New("token expired")
)

// Store Token 存储接口
type Store interface {
	// Save 保存 Token 并设置有效期，单位秒
	Save(token string, expiresIn int64) error
	// Get 获取 Token 内容
	// Token 没有应返回 ErrTokenNotFound
	// Token 过期应返回 ErrTokenExpired
	Get() (string, error)
}

// MemoryStore 内存存储
type MemoryStore struct {
	sync.Map
}

type storeValue struct {
	value string
	timer *time.Timer
}

// Save 保存 Token 并设置有效期，单位秒
func (s *MemoryStore) Save(accessToken string, expiresIn int64) error {
	key := "access_token"
	timer := time.AfterFunc(time.Second*time.Duration(expiresIn), func() {
		s.Delete(key)
	})
	val := storeValue{accessToken, timer}
	s.Store(key, val)
	return nil
}

// Get 获取 Token 内容
func (s *MemoryStore) Get() (string, error) {
	key := "access_token"
	val, found := s.Load(key)
	if !found {
		return "", ErrTokenNotFound
	}
	return val.(storeValue).value, nil
}
