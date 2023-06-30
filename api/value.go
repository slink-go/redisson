package api

import (
	"math"
	"strconv"
	"strings"
)

type Value interface {
	IsEmpty() bool
	String() string
	V() string
	AsInt() int
	AsFloat() float64
	AsBool() bool
}

func NewValue(value string) Value {
	return &redisValue{
		value: value,
	}
}

type redisValue struct {
	value string
}

func (v *redisValue) IsEmpty() bool {
	return strings.TrimSpace(v.value) == ""
}
func (v *redisValue) V() string {
	return v.value
}
func (v *redisValue) String() string {
	return v.value
}
func (v *redisValue) AsInt() int {
	i, err := strconv.ParseInt(v.value, 10, 64)
	if err != nil {
		return math.MinInt
	}
	return int(i)
}
func (v *redisValue) AsFloat() float64 {
	f, err := strconv.ParseFloat(v.value, 64)
	if err != nil {
		return float64(math.MinInt)
	}
	return f
}
func (v *redisValue) AsBool() bool {
	b, err := strconv.ParseBool(v.value)
	if err != nil {
		return false
	}
	return b
}
