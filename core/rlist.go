package core

import (
	"github.com/mediocregopher/radix/v4"
	"go.slink.ws/redisson/api"
	"reflect"
)

func NewRList(key string, client api.Redis) api.RList {
	return &rlist{
		client: client,
		key:    key,
	}
}

type rlist struct {
	client api.Redis
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
func (l *rlist) LPop() (api.Value, error) {
	var value string
	err := l.client.Do(radix.Cmd(&value, "LPOP", l.key))
	return NewValue(value), err
}
func (l *rlist) RPush(items ...any) error {
	return l.client.Do(radix.Cmd(nil, "RPUSH", l.client.AnyArgs(l.key, items...)...))
}
func (l *rlist) RPop() (api.Value, error) {
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
