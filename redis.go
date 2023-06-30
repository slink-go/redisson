package redisson

import (
	"context"
	"errors"
	"fmt"
	"github.com/mediocregopher/radix/v4"
	"time"
)

var ErrRedisClientNotInitialized = errors.New("redis client is not initialized")

type Redis interface {
	Logger

	Close() error

	// helpers

	AnyArgs(key string, args ...any) []string
	StrArgs(key string, args ...string) []string
	Do(cmd radix.Action) error

	// common

	Del(keys ...string) (int, error)
	Expire(key string, ttl time.Duration) (int, error)
	Exists(key ...string) bool
	Keys(filter string) []string
	Touch(keys ...string)
	Type(key string) string

	// basic

	Set(key string, value any) error
	Get(key string) (Value, error)
	Incr(key string) (int, error)
	Decr(key string) (int, error)

	// pub / sub

	PubSub() (radix.PubSubConn, error)

	// wrappers

	RMap(key string) RMap
	RCacheMap(key string) (RCacheMap, error)
	RList(key string) RList
	RSet(key string) RSet
}

type redis struct {
	single   radix.Client
	sentinel *radix.Sentinel
	cluster  *radix.Cluster
	logger   Logger
}

// region - redis

// region - connection

func (r *redis) Close() error {
	var err error
	if r.single != nil {
		err = r.single.Close()
	} else if r.sentinel != nil {
		err = r.sentinel.Close()
	} else if r.cluster != nil {
		err = r.cluster.Close()
	} else {
		err = ErrRedisClientNotInitialized
	}
	return err
}

// endregion
// region - common

func (r *redis) Del(keys ...string) (int, error) {
	var amount int
	var err = r.Do(radix.Cmd(&amount, "DEL", keys...))
	return amount, err
}
func (r *redis) Expire(key string, ttl time.Duration) (int, error) {
	var amount int
	var err = r.Do(radix.Cmd(&amount, "EXPIRE", key, fmt.Sprintf("%0.f", ttl.Seconds())))
	return amount, err
}
func (r *redis) Exists(keys ...string) bool {
	var amount int
	_ = r.Do(radix.Cmd(&amount, "EXISTS", keys...))
	return amount == len(keys) && len(keys) > 0
}
func (r *redis) Keys(filter string) []string {
	if filter == "" {
		filter = "*"
	}
	var keys []string
	_ = r.Do(radix.Cmd(&keys, "KEYS", filter))
	return keys
}
func (r *redis) Touch(keys ...string) {
	_ = r.Do(radix.Cmd(nil, "TOUCH", keys...))
}
func (r *redis) Type(key string) string {
	var value string
	_ = r.Do(radix.Cmd(&value, "TYPE", key))
	return value
}

// endregion
// region - simple

func (r *redis) Set(key string, value any) error {
	return r.Do(radix.Cmd(nil, "SET", key, fmt.Sprintf("%v", value)))
}
func (r *redis) Get(key string) (Value, error) {
	var data string
	var err = r.Do(radix.Cmd(&data, "GET", key))
	return &redisValue{
		value: data,
	}, err
}
func (r *redis) Incr(key string) (int, error) {
	var data int
	var err = r.Do(radix.Cmd(&data, "INCR", key))
	return data, err
}
func (r *redis) Decr(key string) (int, error) {
	var data int
	var err = r.Do(radix.Cmd(&data, "DECR", key))
	return data, err
}

// endregion
// region - pub / sub

func (r *redis) PubSub() (conn radix.PubSubConn, err error) {
	if r.single != nil {
		conn, err = (radix.PersistentPubSubConnConfig{}).New(r.defaultContext(), func() (string, string, error) {
			return r.single.Addr().Network(), r.single.Addr().String(), nil
		})
	} else if r.cluster != nil {
		conn, err = (radix.PersistentPubSubConnConfig{}).New(r.defaultContext(), func() (string, string, error) {
			clients, err := r.cluster.Clients()
			if err != nil {
				return "", "", err
			}
			for addr := range clients {
				return "tcp", addr, nil
			}
			return "", "", errors.New("no clients in the cluster")
		})
		if err != nil {
			return nil, err
		}
	} else if r.sentinel != nil {
		conn, err = (radix.PersistentPubSubConnConfig{}).New(r.defaultContext(), func() (string, string, error) {
			clients, err := r.sentinel.Clients()
			if err != nil {
				return "", "", err
			}
			for addr := range clients {
				return "tcp", addr, nil
			}
			return "", "", errors.New("no clients in the sentinel group")
		})
		if err != nil {
			return nil, err
		}
	}
	return
}

// endregion
// region - wrappers

func (r *redis) RMap(key string) RMap {
	return NewRMap(key, r)
}
func (r *redis) RCacheMap(key string) (RCacheMap, error) {
	return NewRCacheMap(key, r)
}
func (r *redis) RList(key string) RList {
	return NewRList(key, r)
}
func (r *redis) RSet(key string) RSet {
	return NewRSet(key, r)
}

// endregion
// region - helpers

func (r *redis) AnyArgs(key string, args ...any) []string {
	var result []string
	result = append(result, key)
	for _, a := range args {
		switch a.(type) {
		case string:
			result = append(result, a.(string))
		default:
			result = append(result, fmt.Sprintf("%v", a))
		}
	}
	return result
}
func (r *redis) StrArgs(key string, args ...string) []string {
	var result []string
	result = append(result, key)
	result = append(result, args...)
	return result
}
func (r *redis) Do(cmd radix.Action) error {
	var err error
	if r.single != nil {
		err = r.single.Do(r.defaultContext(), cmd)
	} else if r.sentinel != nil {
		err = r.sentinel.Do(r.defaultContext(), cmd)
	} else if r.cluster != nil {
		err = r.cluster.Do(context.Background(), cmd)
	} else {
		err = ErrRedisClientNotInitialized
	}
	return err
}

func (r *redis) defaultContext() context.Context {
	return context.Background()
}

// endregion

// endregion
