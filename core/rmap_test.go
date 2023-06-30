package core

import (
	"fmt"
	"github.com/slink-go/redisson/api"
	"testing"
)

func TestRMap(t *testing.T) {

	r, err := createClient()
	if err != nil {
		t.Error(err)
	}
	defer func(r api.Redis) {
		_ = r.Close()
	}(r)

	m := NewRMap("TEST_MAP", r)
	if m == nil {
		t.Errorf("expected non-null value")
	}

	tp := r.Type("TEST_MAP")
	if tp != "none" {
		t.Errorf("expected 'none' type, received '%s'", tp)
	}

	value, ok := m.Get("key")
	if !ok {
		t.Errorf("expected no value, received '%v'", value)
	}

	err = m.Set("key", "value")
	if err != nil {
		t.Error(err)
	}

	tp = r.Type("TEST_MAP")
	if tp != "hash" {
		t.Errorf("expected 'hash' type, received '%s'", tp)
	}

	value, ok = m.Get("key")
	if !ok || value.IsEmpty() || value.AsString() != "value" {
		t.Errorf("expected 'value', received '%v'", value.String())
	}

	err = m.Del("key")
	if err != nil {
		t.Error(err)
	}

	tp = r.Type("TEST_MAP")
	if tp != "none" {
		t.Errorf("expected 'none' type, received '%s'", tp)
	}

	value, ok = m.Get("key")
	if !ok {
		t.Errorf("expected no value, received '%v'", value)
	}

}
func TestRMapKeysEntries(t *testing.T) {

	r, err := createClient()
	if err != nil {
		t.Error(err)
	}
	defer func(r api.Redis) {
		_ = r.Close()
	}(r)

	m := NewRMap("TEST_MAP", r)
	if m == nil {
		t.Errorf("expected non-null value")
	}

	var testData = make(map[string]any)
	testData["k1"] = "v1"
	testData["k2"] = 0
	testData["k3"] = false
	testData["k4"] = 3.1415

	for k, v := range testData {
		_ = m.Set(k, v)
	}

	if len(m.Keys()) != 4 {
		t.Errorf("expected 4 items, received %d", len(m.Keys()))
	}
	for _, k := range m.Keys() {
		v, ok := testData[k]
		if !ok {
			t.Errorf("expected existing key '%s'", k)
		}
		w, ok := m.Get(k)
		if !ok {
			t.Errorf("expected existing key '%s'", k)
		}
		if w.String() != fmt.Sprintf("%v", v) {
			t.Errorf("expected '%v', received '%v'", v, w)
		}
	}

	if len(m.Entries()) != 4 {
		t.Errorf("expected 4 items, received %d", len(m.Entries()))
	}
	for _, e := range m.Entries() {
		v, ok := testData[e.Key]
		if !ok {
			t.Errorf("found unexpected item %v=%v", e.Key, e.Value)
		}
		if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", e.Value) {
			t.Errorf("expected '%s', received '%s'", v, e.Value)
		}
	}

	_, _ = r.Del("TEST_MAP")

	m = NewRMap("TEST_MAP_2", r)
	if m.Keys() == nil {
		t.Error("expected non-null value")
	}
	if len(m.Keys()) != 0 {
		t.Error("expected empty list")
	}
}
func TestRCacheMap(t *testing.T) {

	r, err := createClient()
	if err != nil {
		t.Error(err)
	}
	defer func(r api.Redis) {
		_ = r.Close()
	}(r)

	m, err := NewRCacheMap("TEST_CACHE_MAP", r)
	if err != nil {
		t.Error(err)
	}
	if m == nil {
		t.Errorf("expected non-null value")
	}

	err = m.Set("key1", "value1")
	if err != nil {
		t.Error(err)
	}

	value, ok := m.Get("key1")
	if !ok {
		t.Errorf("expected existing key '%s'", "key1")
	}
	if value.IsEmpty() {
		t.Errorf("expected non-empty value for '%s'", "key1")
	}

	keys := m.Keys()
	if keys == nil {
		t.Errorf("expected non-null keys list")
	}
	if len(keys) == 0 {
		t.Errorf("expected non-empty keys list")
	}
	if len(keys) != 1 {
		t.Errorf("expected one item only")
	}
	if keys[0] != "key1" {
		t.Errorf("expected '%s', received '%s'", "key1", keys[0])
	}

	entries := m.Entries()
	if entries == nil {
		t.Errorf("expected non-null entries list")
	}
	if len(entries) == 0 {
		t.Errorf("expected non-empty entries list")
	}
	if len(entries) != 1 {
		t.Errorf("expected one item only")
	}
	if entries[0].Key != "key1" {
		t.Errorf("expected '%s', received '%s'", "key1", keys[0])
	}
	if entries[0].Value.IsEmpty() {
		t.Errorf("expected non-empty value")
	}
	if entries[0].Value.AsString() != "value1" {
		t.Errorf("expected '%s', received '%s'", "value1", entries[0].Value)
	}

	_ = m.Del("key1")
	keys = m.Keys()
	if keys == nil {
		t.Errorf("expected non-null keys list")
	}
	if len(keys) != 0 {
		t.Errorf("expected empty keys list")
	}

	m.Destroy()

}
