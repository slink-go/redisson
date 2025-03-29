package core

import (
	"github.com/mediocregopher/radix/v4"
	"go.slink.ws/redisson/api"
)

type rset struct {
	client api.Redis
	key    string
}

func NewRSet(key string, client api.Redis) api.RSet {
	return &rset{
		client: client,
		key:    key,
	}
}

func (s *rset) Size() int {
	var result int
	_ = s.client.Do(radix.Cmd(&result, "SCARD", s.key))
	return result
}
func (s *rset) Add(values ...any) error {
	return s.client.Do(radix.Cmd(nil, "SADD", s.client.AnyArgs(s.key, values...)...))
}
func (s *rset) Has(value any) bool {
	var result int
	_ = s.client.Do(radix.Cmd(&result, "SISMEMBER", s.client.AnyArgs(s.key, value)...))
	return result > 0
}
func (s *rset) Del(values ...any) error {
	return s.client.Do(radix.Cmd(nil, "SREM", s.client.AnyArgs(s.key, values...)...))
}
func (s *rset) Items() []api.Value {
	var result []string
	_ = s.client.Do(radix.Cmd(&result, "SMEMBERS", s.key))
	var values []api.Value
	for _, v := range result {
		values = append(values, NewValue(v))
	}
	return values
}
