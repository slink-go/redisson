package core

import (
	"github.com/slink-go/redisson/api"
	"testing"
)

func TestRSet(t *testing.T) {

	r, err := createClient()
	if err != nil {
		t.Error(err)
	}
	defer func(r api.Redis) {
		_ = r.Close()
	}(r)

	s := NewRSet("TEST_SET", r)

	if s.Size() != 0 {
		t.Errorf("expected 0, received %d", s.Size())
	}

	testData := make(map[string]struct{})
	testData["true"] = struct{}{}
	testData["1"] = struct{}{}
	testData["two"] = struct{}{}

	_ = s.Add(true, 1, "two")

	if s.Size() != 3 {
		t.Errorf("expected 3, received %d", s.Size())
	}

	has := s.Has(1)
	if !has {
		t.Errorf("expected 'true', received '%v'", has)
	}
	has = s.Has(2)
	if has {
		t.Errorf("expected 'false', received '%v'", has)
	}

	if len(s.Items()) != s.Size() {
		t.Errorf("mismatched cardinality & size: '%v' vs '%v'", len(s.Items()), s.Size())
	}
	for _, item := range s.Items() {
		_, ok := testData[item.AsString()]
		if !ok {
			t.Errorf("unexpected item found '%v'", item)
		}
	}

	_ = s.Del(1, true)
	has = s.Has(1)
	if has {
		t.Errorf("expected 'false', received '%v'", has)
	}

	_, _ = r.Del("TEST_SET")

}
