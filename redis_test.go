package redisson

import (
	"fmt"
	"github.com/stvp/tempredis"
	"os"
	"strconv"
	"testing"
	"time"
)

const testServerPort = 51200
const testServerHost = "127.0.0.1"

var server *tempredis.Server

func TestMain(m *testing.M) {
	var err error
	os.Remove("dump.rdb")
	server, err = tempredis.Start(tempredis.Config{
		"port": strconv.Itoa(testServerPort),
	})
	if err != nil {
		panic(err)
	}
	result := m.Run()
	_ = server.Term()
	os.Exit(result)
}

func createClient() (Redis, error) {
	return NewConfig().
		WithName("TEST-SINGLE-CLIENT").
		WithDb(0).
		WithPoolSize(5).
		NewSingle(fmt.Sprintf("%s:%d", testServerHost, testServerPort))
}

func TestNewClient(t *testing.T) {
	r, err := createClient()
	if err != nil {
		t.Error(err)
	}
	err = r.Close()
	if err != nil {
		t.Error(err)
	}
}
func TestSetGetDelete(t *testing.T) {
	r, err := createClient()
	if err != nil {
		t.Error(err)
	}
	defer func(r Redis) {
		_ = r.Close()
	}(r)
	err = r.Set("TEST_KEY", "TEST_VALUE")
	if err != nil {
		t.Error(err)
	}
	v, err := r.Get("TEST_KEY")
	if err != nil {
		t.Error(err)
	}
	if v.String() != "TEST_VALUE" {
		t.Errorf("expected '%s', but received '%s'", "TEST_VALUE", v.String())
	}
	i, err := r.Del("TEST_KEY")
	if err != nil {
		t.Error(err)
	}
	if i != 1 {
		t.Errorf("expected %d, but received %d", 1, i)
	}
	v, err = r.Get("TEST_KEY")
	if err != nil {
		t.Error(err)
	}
	if v.String() != "" {
		t.Errorf("expected empty value, but received '%s'", v.String())
	}
}
func TestExistsExpireDelete(t *testing.T) {
	r, err := createClient()
	if err != nil {
		t.Error(err)
	}
	defer func(r Redis) {
		_ = r.Close()
	}(r)

	value, err := r.Get("TEST_KEY")
	if err != nil {
		t.Error(err)
	}
	if !value.IsEmpty() {
		t.Errorf("expected empty value, received '%s'", value)
	}

	exists := r.Exists("TEST_KEY")
	if exists {
		t.Errorf("expected non-existent key")
	}

	err = r.Set("TEST_KEY", "TEST_VALUE")
	if err != nil {
		t.Error(err)
	}

	exists = r.Exists("TEST_KEY")
	if !exists {
		t.Errorf("expected existent key")
	}

	value, err = r.Get("TEST_KEY")
	if err != nil {
		t.Error(err)
	}
	if value.IsEmpty() {
		t.Errorf("expected non-empty value")
	}
	if value.AsString() != "TEST_VALUE" {
		t.Errorf("expected '%s', received '%s'", "TEST_VALUE", value)
	}

	_, err = r.Expire("TEST_KEY", 100*time.Millisecond)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(250 * time.Millisecond)

	value, err = r.Get("TEST_KEY")
	if err != nil {
		t.Error(err)
	}
	if !value.IsEmpty() {
		t.Errorf("expected empty value, received '%s'", value)
	}

	exists = r.Exists("TEST_KEY")
	if exists {
		t.Errorf("expected non-existent key")
	}
}
func TestIncrDecrDelete(t *testing.T) {
	r, err := createClient()
	if err != nil {
		t.Error(err)
	}
	defer func(r Redis) {
		_ = r.Close()
	}(r)

	i, err := r.Incr("TEST_KEY")
	if err != nil {
		t.Error(err)
	}
	if i != 1 {
		t.Errorf("expected %d, but received %d", 1, i)
	}

	i, err = r.Incr("TEST_KEY")
	if err != nil {
		t.Error(err)
	}
	if i != 2 {
		t.Errorf("expected %d, but received %d", 2, i)
	}

	i, err = r.Decr("TEST_KEY")
	if err != nil {
		t.Error(err)
	}
	if i != 1 {
		t.Errorf("expected %d, but received %d", 1, i)
	}

	i, err = r.Decr("TEST_KEY")
	if err != nil {
		t.Error(err)
	}
	if i != 0 {
		t.Errorf("expected %d, but received %d", 0, i)
	}

	i, err = r.Del("TEST_KEY")
	if err != nil {
		t.Error(err)
	}
	if i != 1 {
		t.Errorf("expected %d, but received %d", 1, i)
	}

	v, err := r.Get("TEST_KEY")
	if err != nil {
		t.Error(err)
	}
	if v.String() != "" {
		t.Errorf("expected empty value, but received '%s'", v.String())
	}
}
func TestKeys(t *testing.T) {
	r, err := createClient()
	if err != nil {
		t.Error(err)
	}
	defer func(r Redis) {
		_ = r.Close()
	}(r)

	keys := r.Keys("")
	if keys == nil {
		t.Errorf("expected non-null value")
	}
	if len(keys) != 0 {
		t.Errorf("expected empty array")
	}

	r.Set("TEST_KEY_1", "TEST_VALUE_1")
	r.Set("TEST_KEY_2", "TEST_VALUE_2")

	keys = r.Keys("")
	if keys == nil {
		t.Errorf("expected non-null value")
	}
	if len(keys) != 2 {
		t.Errorf("expected two items")
	}

	keys = r.Keys("*_2")
	if keys == nil {
		t.Errorf("expected non-null value")
	}
	if len(keys) != 1 {
		t.Errorf("expected one item")
	}
	if keys[0] != "TEST_KEY_2" {
		t.Errorf("expected '%s', received '%s'", "TEST_KEY_2", keys[0])
	}

	_, _ = r.Del("TEST_KEY_1", "TEST_KEY_2")
}
func TestType(t *testing.T) {
	r, err := createClient()
	if err != nil {
		t.Error(err)
	}
	defer func(r Redis) {
		_ = r.Close()
	}(r)
	r.Set("TEST_KEY", "TEST_VALUE")
	tp := r.Type("TEST_KEY")
	if tp == "" {
		t.Errorf("expected non-empty value")
	}
	if tp != "string" {
		t.Errorf("expected 'string', received '%s'", tp)
	}
	_, _ = r.Del("TEST_KEY")
}
func TestAnyArgs(t *testing.T) {
	r, err := createClient()
	if err != nil {
		t.Error(err)
	}
	defer func(r Redis) {
		_ = r.Close()
	}(r)

	v := r.AnyArgs("key", 1, 3.14, false, "yes")
	if v[0] != "key" {
		t.Errorf("expected 'key', received '%v'", v[0])
	}
	if v[1] != "1" {
		t.Errorf("expected '1', received '%v'", v[1])
	}
	if v[2] != "3.14" {
		t.Errorf("expected '3.14', received '%v'", v[2])
	}
	if v[3] != "false" {
		t.Errorf("expected 'false', received '%v'", v[3])
	}
	if v[4] != "yes" {
		t.Errorf("expected 'yes', received '%v'", v[4])
	}
}
func TestStrArgs(t *testing.T) {
	r, err := createClient()
	if err != nil {
		t.Error(err)
	}
	defer func(r Redis) {
		_ = r.Close()
	}(r)

	v := r.StrArgs("key", "arg1", "arg2")
	if v[0] != "key" {
		t.Errorf("expected 'key', received '%v'", v[0])
	}
	if v[1] != "arg1" {
		t.Errorf("expected 'arg1', received '%v'", v[1])
	}
	if v[2] != "arg2" {
		t.Errorf("expected 'arg2', received '%v'", v[2])
	}
}

func TestRMap(t *testing.T) {

	r, err := createClient()
	if err != nil {
		t.Error(err)
	}
	defer func(r Redis) {
		_ = r.Close()
	}(r)

	m := r.RMap("TEST_MAP")
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
	defer func(r Redis) {
		_ = r.Close()
	}(r)

	m := r.RMap("TEST_MAP")
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

}

func TestListL(t *testing.T) {

	r, err := createClient()
	if err != nil {
		t.Error(err)
	}
	defer func(r Redis) {
		_ = r.Close()
	}(r)

	l := r.RList("TEST_LIST")

	if l.Len() != 0 {
		t.Errorf("expected 0, received %d", l.Len())
	}

	l.LPush(true, 1, "two")

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
func TestListR(t *testing.T) {

	r, err := createClient()
	if err != nil {
		t.Error(err)
	}
	defer func(r Redis) {
		_ = r.Close()
	}(r)

	l := r.RList("TEST_LIST")

	if l.Len() != 0 {
		t.Errorf("expected 0, received %d", l.Len())
	}

	l.RPush(true, 1, "two")

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
func TestListLRO(t *testing.T) {
	r, err := createClient()
	if err != nil {
		t.Error(err)
	}
	defer func(r Redis) {
		_ = r.Close()
	}(r)

	l := r.RList("TEST_LIST")

	if l.Len() != 0 {
		t.Errorf("expected 0, received %d", l.Len())
	}

	arr := []any{true, 1, "two"}
	l.LPushRO(arr...)

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
