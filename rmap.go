package redisson

import (
	"context"
	"errors"
	"fmt"
	"github.com/mediocregopher/radix/v4"
	"sync"
	"time"
)

// region - map entry

type Entry struct {
	Key   string
	Value Value
}

// endregion
// region - RMap

type RMap interface {
	Set(key string, value any) error
	Get(key string) (Value, bool)
	Del(keys ...string) error
	Keys() []string
	Entries() []Entry
}

func NewRMap(key string, client Redis) RMap {
	return &rmap{
		client: client,
		key:    key,
	}
}

type rmap struct {
	client Redis
	key    string
}

func (m *rmap) Set(key string, value any) error {
	return m.client.Do(radix.Cmd(nil, "HSET", m.key, key, fmt.Sprintf("%v", value)))
}
func (m *rmap) Get(key string) (Value, bool) {
	var value string
	var ok = true
	err := m.client.Do(radix.Cmd(&value, "HGET", m.key, key))
	if err != nil {
		m.client.Warning("RMap get error: %s", err.Error())
		ok = false
	}
	return NewValue(value), ok
}
func (m *rmap) Del(keys ...string) error {
	return m.client.Do(radix.Cmd(nil, "HDEL", m.client.StrArgs(m.key, keys...)...))
}
func (m *rmap) Keys() []string {
	var result []string
	_ = m.client.Do(radix.Cmd(&result, "HKEYS", m.key))
	if result == nil {
		result = []string{}
	}
	return result
}
func (m *rmap) Entries() []Entry {
	var result map[string]string
	_ = m.client.Do(radix.Cmd(&result, "HGETALL", m.key))
	var values []Entry
	for k, v := range result {
		values = append(values, Entry{
			Key:   k,
			Value: NewValue(v),
		})
	}
	return values
}

// endregion
// region - RCacheMap

// see: https://redis.io/docs/manual/keyspace-notifications/

const keySpaceTopicFormat = "__keyspace@*__:%s"

type syncState uint8

const (
	syncNeeded syncState = iota
	syncPending
	syncInProgress
	syncComplete
)

type RCacheMap interface {
	RMap
	Destroy()
}

type rcachemap struct {
	client    Redis
	key       string
	rwMutex   sync.RWMutex
	syncMutex sync.RWMutex
	syncState syncState
	cache     map[string]Value
	redisChn  chan radix.PubSubMessage
	doneChn   chan *struct{}
	psconn    radix.PubSubConn
}

func NewRCacheMap(key string, client Redis) (RCacheMap, error) {
	m := &rcachemap{
		client:    client,
		key:       key,
		syncState: syncNeeded,
		cache:     make(map[string]Value),
		redisChn:  make(chan radix.PubSubMessage, 1),
		doneChn:   make(chan *struct{}, 1),
	}
	err := m.run()
	if err != nil {
		m.client.Error("could not start subscription for key %s", key)
		m.Destroy()
	}
	return m, err
}

func (m *rcachemap) Set(key string, value any) error {
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
	m.client.Debug("+ set start: %s", key)
	err := m.client.Do(radix.Cmd(nil, "HSET", m.key, key, fmt.Sprintf("%v", value)))
	m.syncState = syncNeeded
	m.client.Debug("+ set end: %s %d", key, m.syncState)
	return err
}
func (m *rcachemap) Get(key string) (Value, bool) {
	m.wait()
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	v, ok := m.cache[key]
	return v, ok
}
func (m *rcachemap) Del(keys ...string) error {
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
	m.client.Debug("- del start: %s", keys)
	err := m.client.Do(radix.Cmd(nil, "HDEL", m.client.StrArgs(m.key, keys...)...))
	m.syncState = syncNeeded
	m.client.Debug("- del end: %s", keys)
	return err
}
func (m *rcachemap) Keys() []string {
	m.wait()
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	var result []string
	for k := range m.cache {
		result = append(result, k)
	}
	if result == nil {
		result = []string{}
	}
	return result
}
func (m *rcachemap) Entries() []Entry {
	m.wait()
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	var result []Entry
	for k, v := range m.cache {
		result = append(result, Entry{
			Key:   k,
			Value: v,
		})
	}
	return result
}
func (m *rcachemap) Destroy() {
	m.doneChn <- &struct{}{}
	if m.psconn != nil {
		_ = m.psconn.PUnsubscribe(context.Background())
		_ = m.psconn.Close()
	}
}

func (m *rcachemap) run() error {
	m.client.Debug("starting RCacheMap background process for %s", m.key)
	var err error
	err = m.client.Do(radix.Cmd(nil, "config", "set", "notify-keyspace-events", "KEAn"))
	if err != nil {
		return err
	}
	m.psconn, err = m.client.PubSub()
	if err != nil {
		return err
	}
	err = m.psconn.PSubscribe(context.Background(), fmt.Sprintf(keySpaceTopicFormat, m.key))
	if err != nil {
		return err
	}
	go func() {
		timer := time.NewTimer(100 * time.Millisecond)
		for {
			select {
			case v := <-m.doneChn:
				if v != nil {
					timer.Stop()
					m.client.Debug("stopping RCacheMap background process for %s", m.key)
					close(m.doneChn)
					close(m.redisChn)
				}
				break
			case msg := <-m.redisChn:
				if msg.Channel != "" {
					m.handleMessage(msg)
				}
			case <-timer.C:
				if m.syncState == syncPending {
					m.sync()
				}
			default:
				m.checkSubscription()
				if m.syncState == syncNeeded {
					m.client.Debug("reschedule sync")
					timer.Stop()
					timer.Reset(100 * time.Millisecond)
					m.syncState = syncPending
				}
			}
		}
	}()
	m.sync()
	return err
}
func (m *rcachemap) sync() {
	m.client.Debug("sync start")
	m.syncMutex.Lock()
	defer m.syncMutex.Unlock()
	m.syncState = syncInProgress
	var keys []string
	err := m.client.Do(radix.Cmd(&keys, "HKEYS", m.key))
	if err != nil {
		m.client.Warning("sync keys error: %s", err.Error())
	} else {
		m.cache = make(map[string]Value)
		for _, key := range keys {
			var value string
			err := m.client.Do(radix.Cmd(&value, "HGET", m.key, key))
			if err != nil {
				m.client.Warning("sync error for key %s: %s", key, err.Error())
				continue
			}
			m.client.Debug("cache map %s: sync key %s=%v", m.key, key, value)
			m.cache[key] = NewValue(value)
		}
	}
	m.client.Debug("sync end")
	m.syncState = syncComplete
}
func (m *rcachemap) checkSubscription() {
	//m.client.Warning("subscription check")
	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()
	msg, err := m.psconn.Next(ctx)
	if errors.Is(err, context.DeadlineExceeded) {
		//m.client.Debug("subscription check timeout")
	} else if err != nil {
		m.client.Warning("subscription check error: %s", err.Error())
	} else {
		m.redisChn <- msg
	}
}
func (m *rcachemap) handleMessage(msg radix.PubSubMessage) {
	message := string(msg.Message)
	switch message {
	case "new":
		fallthrough
	case "set":
		fallthrough
	case "del":
		fallthrough
	case "hset":
		fallthrough
	case "hdel":
		m.client.Debug("handle: %s %s", msg.Channel, message)
		m.syncState = syncNeeded
	default:
		m.client.Debug("skip: %s %s", msg.Channel, message)
	}
}
func (m *rcachemap) wait() {
	if m.syncState == syncComplete {
		return
	}
	for {
		select {
		case <-time.After(5 * time.Second):
			m.client.Notice("sync wait timeout")
			m.syncState = syncComplete
			return
		default:
			if m.syncState == syncComplete {
				m.client.Debug("wait: sync complete")
				return
			}
		}
	}
}

//endregion
