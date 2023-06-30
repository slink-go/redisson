package redisson

import (
	"fmt"
	"math"
	"strconv"
	"strings"
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

func NewValue(value any) Value {
	return &redisValue{
		value: value,
	}
}

type redisValue struct {
	value any
}

func (v *redisValue) IsEmpty() bool {
	return strings.TrimSpace(v.String()) == ""
}
func (v *redisValue) V() any {
	return v.value
}
func (v *redisValue) String() string {
	return fmt.Sprintf("%v", v.value)
}
func (v *redisValue) AsString() string {
	return v.String()
}
func (v *redisValue) AsInt() int {
	i, err := strconv.ParseInt(v.String(), 10, 64)
	if err != nil {
		return math.MinInt
	}
	return int(i)
}
func (v *redisValue) AsFloat() float64 {
	f, err := strconv.ParseFloat(v.String(), 64)
	if err != nil {
		return float64(math.MinInt)
	}
	return f
}
func (v *redisValue) AsBool() bool {
	b, err := strconv.ParseBool(v.String())
	if err != nil {
		return false
	}
	return b
}
