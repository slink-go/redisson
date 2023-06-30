package redisson

import (
	"github.com/mediocregopher/radix/v4"
	"reflect"
)

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

func NewRList(key string, client Redis) RList {
	return &rlist{
		client: client,
		key:    key,
	}
}

type rlist struct {
	client Redis
	key    string
}

func (l *rlist) Len() int {
	var value int
	_ = l.client.Do(radix.Cmd(&value, "LLEN", l.key))
	return value
}
func (l *rlist) LPush(items ...any) error {
	return l.client.Do(radix.Cmd(nil, "LPUSH", l.client.AnyArgs(l.key, items...)...))
}
func (l *rlist) LPushRO(items ...any) error {
	i2 := make([]any, len(items))
	copy(i2, items)
	ReverseSlice(i2)
	return l.client.Do(radix.Cmd(nil, "LPUSH", l.client.AnyArgs(l.key, i2...)...))
}
func (l *rlist) LPop() (Value, error) {
	var value string
	err := l.client.Do(radix.Cmd(&value, "LPOP", l.key))
	return NewValue(value), err
}
func (l *rlist) RPush(items ...any) error {
	return l.client.Do(radix.Cmd(nil, "RPUSH", l.client.AnyArgs(l.key, items...)...))
}
func (l *rlist) RPop() (Value, error) {
	var value string
	err := l.client.Do(radix.Cmd(&value, "RPOP", l.key))
	return NewValue(value), err
}

func ReverseSlice(s interface{}) {
	size := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}
