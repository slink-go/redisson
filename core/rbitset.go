package core

import (
	"fmt"
	"github.com/mediocregopher/radix/v4"
	"go.slink.ws/redisson/api"
	"strings"
)

type rbitset struct {
	client api.Redis
	key    string
}

func NewRBitSet(key string, client api.Redis) api.RBitSet {
	return &rbitset{
		client: client,
		key:    key,
	}
}

func (bs *rbitset) Set(idx uint32, value any) (bool, error) {
	var result int
	err := bs.client.Do(radix.Cmd(&result, "SETBIT", bs.client.AnyArgs(bs.key, idx, value)...))
	return result > 0, err
}
func (bs *rbitset) Get(idx uint32) (bool, error) {
	var result int
	err := bs.client.Do(radix.Cmd(&result, "GETBIT", bs.key, fmt.Sprintf("%v", idx)))
	return result > 0, err
}
func (bs *rbitset) BitCount() int {
	var result int
	_ = bs.client.Do(radix.Cmd(&result, "BITCOUNT", bs.key))
	return result
}
func (bs *rbitset) BitCountRange(start, end int, unit string) (int, error) {
	if start > end && end >= 0 {
		v := end
		end = start
		start = v
	}
	if unit == "" {
		unit = "BYTE"
	}
	switch strings.ToUpper(unit) {
	case "BIT":
	case "BYTE":
	default:
		return 0, fmt.Errorf("invalid unit '%s', supported values are: BIT, BYTE", unit)
	}
	var result int
	err := bs.client.Do(radix.Cmd(&result, "BITCOUNT", bs.client.AnyArgs(bs.key, start, end, unit)...))
	return result, err
}
