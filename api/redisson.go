package api

import (
	"github.com/mediocregopher/radix/v4"
	"time"
)

type Value interface {
	IsEmpty() bool
	String() string
	V() any
	AsString() string
	AsInt() int
	AsFloat() float64
	AsBool() bool
}

type MapEntry struct {
	Key   string
	Value Value
}

type RList interface {

	// Len returns length of a list
	Len() int

	// LPush adds items to list tail in given order
	LPush(items ...any) error

	// LPushRO adds items to list tail in reversed order
	//       i.e. first item in passed list will be added last
	LPushRO(items ...any) error

	// LPop get item from list tail
	LPop() (Value, error)

	// RPush adds items to list head in given order
	RPush(items ...any) error

	// RPop get item from list head
	RPop() (Value, error)
}
type RSet interface {
	Size() int
	Add(value ...any) error
	Has(value any) bool
	Del(keys ...any) error
	Items() []Value
}
type RBitSet interface {
	Set(idx uint32, value any) (bool, error)
	Get(idx uint32) (bool, error)
	BitCount() int
	BitCountRange(start, end int, unit string) (int, error)
}
type RMap interface {
	Set(key string, value any) error
	Get(key string) (Value, bool)
	Del(keys ...string) error
	Keys() []string
	Entries() []MapEntry
}
type RCacheMap interface {
	RMap
	Destroy()
}

type Redis interface {
	Logger

	EnableKeyEventNotifications() error
	EnableKeyEventNotificationsOfTypes(types string) error
	DisableKeyEventNotifications() error

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
}
