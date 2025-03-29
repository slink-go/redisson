package core

import (
	"go.slink.ws/redisson/api"
	"testing"
)

func TestRBitSet(t *testing.T) {

	r, err := createClient()
	if err != nil {
		t.Error(err)
	}
	defer func(r api.Redis) {
		_ = r.Close()
	}(r)

	bs := NewRBitSet("TEST_BITSET", r)

	if bs.BitCount() != 0 {
		t.Errorf("expected 0, received %d", bs.BitCount())
	}

	v, err := bs.Set(1, 1)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if v {
		t.Errorf("expected 'false', received 'true'")
	}

	v, err = bs.Set(5, 1)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if v {
		t.Errorf("expected 'false', received 'true'")
	}

	v, err = bs.Set(7, 1)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if v {
		t.Errorf("expected 'false', received 'true'")
	}

	v, err = bs.Set(12, 1)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if v {
		t.Errorf("expected 'false', received 'true'")
	}

	v, _ = bs.Get(0)
	if v != false {
		t.Errorf("expected 'true', received 'false'")
	}

	v, _ = bs.Get(1)
	if v != true {
		t.Errorf("expected 'true', received %v", v)
	}

	v, _ = bs.Get(2)
	if v != false {
		t.Errorf("expected 'false', received %v", v)
	}

	v, _ = bs.Get(5)
	if v != true {
		t.Errorf("expected 'true', received %v", v)
	}

	v, _ = bs.Get(12)
	if v != true {
		t.Errorf("expected 'true', received %v", v)
	}

	c := bs.BitCount()
	if c != 4 {
		t.Errorf("expected 4, received %v", c)
	}

	c, err = bs.BitCountRange(0, 1, "bit")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if c != 1 {
		t.Errorf("expected 1, received %v", c)
	}

	c, err = bs.BitCountRange(2, 0, "BIT")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if c != 1 {
		t.Errorf("expected 1, received %v", c)
	}

	c, err = bs.BitCountRange(0, 5, "BIT")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if c != 2 {
		t.Errorf("expected 2, received %v", c)
	}

	c, err = bs.BitCountRange(0, 1, "BITE")
	if err == nil {
		t.Errorf("expected unsupported unit error")
	}

	c, err = bs.BitCountRange(0, 0, "BYTE")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if c != 3 {
		t.Errorf("expected 3, received %v", c)
	}

	c, err = bs.BitCountRange(0, 1, "")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if c != 4 {
		t.Errorf("expected 4, received %v", c)
	}

	_, _ = r.Del("TEST_BITSET")

}
