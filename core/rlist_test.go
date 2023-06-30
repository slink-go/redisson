package core

import (
	"github.com/slink-go/redisson/api"
	"testing"
)

func TestRListL(t *testing.T) {

	r, err := createClient()
	if err != nil {
		t.Error(err)
	}
	defer func(r api.Redis) {
		_ = r.Close()
	}(r)

	l := NewRList("TEST_LIST", r)

	if l.Len() != 0 {
		t.Errorf("expected 0, received %d", l.Len())
	}

	_ = l.LPush(true, 1, "two")

	if l.Len() != 3 {
		t.Errorf("expected 3, received %d", l.Len())
	}

	v, err := l.RPop()
	if err != nil {
		t.Error(err)
	}
	if v.AsBool() != true {
		t.Errorf("expected 'true', received '%v'", v)
	}

	v, err = l.RPop()
	if err != nil {
		t.Error(err)
	}
	if v.AsInt() != 1 {
		t.Errorf("expected '1', received '%v'", v)
	}

	v, err = l.RPop()
	if err != nil {
		t.Error(err)
	}
	if v.AsString() != "two" {
		t.Errorf("expected 'two', received '%v'", v)
	}

	_, _ = r.Del("TEST_LIST")

}
func TestRListR(t *testing.T) {

	r, err := createClient()
	if err != nil {
		t.Error(err)
	}
	defer func(r api.Redis) {
		_ = r.Close()
	}(r)

	l := NewRList("TEST_LIST", r)

	if l.Len() != 0 {
		t.Errorf("expected 0, received %d", l.Len())
	}

	_ = l.RPush(true, 1, "two")

	if l.Len() != 3 {
		t.Errorf("expected 3, received %d", l.Len())
	}

	v, err := l.LPop()
	if err != nil {
		t.Error(err)
	}
	if v.AsBool() != true {
		t.Errorf("expected 'true', received '%v'", v)
	}

	v, err = l.LPop()
	if err != nil {
		t.Error(err)
	}
	if v.AsInt() != 1 {
		t.Errorf("expected '1', received '%v'", v)
	}

	v, err = l.LPop()
	if err != nil {
		t.Error(err)
	}
	if v.AsString() != "two" {
		t.Errorf("expected 'two', received '%v'", v)
	}

	_, _ = r.Del("TEST_LIST")

}
func TestRListLRO(t *testing.T) {
	r, err := createClient()
	if err != nil {
		t.Error(err)
	}
	defer func(r api.Redis) {
		_ = r.Close()
	}(r)

	l := NewRList("TEST_LIST", r)

	if l.Len() != 0 {
		t.Errorf("expected 0, received %d", l.Len())
	}

	arr := []any{true, 1, "two"}
	_ = l.LPushRO(arr...)

	if l.Len() != 3 {
		t.Errorf("expected 3, received %d", l.Len())
	}

	v, err := l.LPop()
	if err != nil {
		t.Error(err)
	}
	if v.AsBool() != true {
		t.Errorf("expected 'true', received '%v'", v)
	}

	v, err = l.LPop()
	if err != nil {
		t.Error(err)
	}
	if v.AsInt() != 1 {
		t.Errorf("expected '1', received '%v'", v)
	}

	v, err = l.LPop()
	if err != nil {
		t.Error(err)
	}
	if v.AsString() != "two" {
		t.Errorf("expected 'two', received '%v'", v)
	}

	_, _ = r.Del("TEST_LIST")

}
