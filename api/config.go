package api

import (
	"context"
	"fmt"
	"github.com/mediocregopher/radix/v4"
	"time"
)

// region - config

const defaultPoolSize = 5
const defaultPingInterval = 5 * time.Second

type config struct {
	name         string
	db           int
	poolSize     int
	pingInterval time.Duration
	user         string
	password     string
	logger       Logger
}

func NewConfig() *config {
	return &config{
		poolSize:     defaultPoolSize,
		pingInterval: defaultPingInterval,
	}
}
func (c *config) WithLogger(logger Logger) *config {
	c.logger = logger
	return c
}
func (c *config) WithDb(db int) *config {
	c.db = db
	return c
}
func (c *config) WithName(name string) *config {
	c.name = name
	return c
}
func (c *config) WithPoolSize(sz int) *config {
	c.poolSize = sz
	return c
}
func (c *config) WithPingInterval(interval time.Duration) *config {
	c.pingInterval = interval
	return c
}
func (c *config) WithAuth(user, password string) *config {
	c.user = user
	c.password = password
	return c
}

func (c *config) NewSingle(addr string) (Redis, error) {
	client, err := (radix.PoolConfig{
		Size:         c.poolSize,
		PingInterval: c.pingInterval,
		Dialer: radix.Dialer{
			CustomConn: c.customConn,
			AuthUser:   c.user,
			AuthPass:   c.password,
			SelectDB:   fmt.Sprintf("%d", c.db),
		},
	}).New(context.Background(), "tcp", addr)
	if err != nil {
		return nil, err
	}
	return &redis{
		single: client,
		logger: c.logger,
	}, nil
}
func (c *config) NewCluster(addr ...string) (Redis, error) {
	cfg := radix.ClusterConfig{
		PoolConfig: radix.PoolConfig{
			Size:         c.poolSize,
			PingInterval: c.pingInterval,
			Dialer: radix.Dialer{
				CustomConn: c.customConn,
				AuthUser:   c.user,
				AuthPass:   c.password,
				SelectDB:   fmt.Sprintf("%d", c.db),
			},
		},
	}
	client, err := cfg.New(context.Background(), addr)
	if err != nil {
		return nil, err
	}
	return &redis{
		cluster: client,
		logger:  c.logger,
	}, err
}
func (c *config) NewSentinel(name string, addr ...string) (Redis, error) {
	cfg := radix.SentinelConfig{
		PoolConfig: radix.PoolConfig{
			Size:         c.poolSize,
			PingInterval: c.pingInterval,
			Dialer: radix.Dialer{
				CustomConn: c.customConn,
				AuthUser:   c.user,
				AuthPass:   c.password,
				SelectDB:   fmt.Sprintf("%d", c.db),
			},
		},
	}
	client, err := cfg.New(context.Background(), name, addr)
	if err != nil {
		return nil, err
	}
	return &redis{
		sentinel: client,
		logger:   c.logger,
	}, err
}

func (c *config) customConn(ctx context.Context, network, addr string) (radix.Conn, error) {
	cl, err := radix.Dial(ctx, network, addr)
	if err != nil {
		return nil, err
	}
	if c.name != "" {
		err = cl.Do(ctx, radix.Cmd(nil, "CLIENT", "SETNAME", c.name))
	}
	return cl, err
}

// endregion
