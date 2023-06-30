package stack

import (
	"github.com/slink-go/redisson/api"
)

// https://redis.io/commands/?name=bf.

type RBloomFilter interface {
}

type rBloomFilter struct {
	key    string
	client api.Redis
}

func NewRBloomFilter(key string, client api.Redis) RBloomFilter {
	return &rBloomFilter{
		key:    key,
		client: client,
	}
}
