package redisson

import (
	"testing"
	"time"
)

const testName = "TEST_NAME"
const testDb = 1
const testPoolSize = 10
const testPingInterval = time.Second
const testUser = "TEST_USER"
const testPass = "TEST_PASS"

func TestConfig(t *testing.T) {

	cfg := NewConfig()
	if cfg == nil {
		t.Errorf("expected non-null value, received '%v'", cfg)
	}

	cfg.WithName(testName)
	if cfg.name != testName {
		t.Errorf("expected '%s', received '%s'", testName, cfg.name)
	}

	cfg.WithDb(testDb)
	if cfg.db != testDb {
		t.Errorf("expected '%d', received '%d'", testDb, cfg.db)
	}

	cfg.WithPoolSize(testPoolSize)
	if cfg.poolSize != testPoolSize {
		t.Errorf("expected '%d', received '%d'", testPoolSize, cfg.poolSize)
	}

	cfg.WithPingInterval(testPingInterval)
	if cfg.pingInterval != testPingInterval {
		t.Errorf("expected '%d', received '%d'", testPingInterval, cfg.pingInterval)
	}

	cfg.WithAuth(testUser, testPass)
	if cfg.user != testUser {
		t.Errorf("expected '%s', received '%s'", testUser, cfg.user)
	}
	if cfg.password != testPass {
		t.Errorf("expected '%s', received '%s'", testPass, cfg.password)
	}

}
