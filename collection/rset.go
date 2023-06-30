package collection

import "github.com/slink-go/redisson/api"

type RSet interface {
	Add(value any) error
	Has(value any) (bool, error)
	Del(keys ...any) error
	Items() []api.Value
}

type rset struct {
	client api.Redis
	key    string
}

func NewRSet(key string, client api.Redis) RSet {
	return &rset{
		client: client,
		key:    key,
	}
}

func (s *rset) Add(value any) error {
	return nil
}
func (s *rset) Has(value any) (bool, error) {
	return false, nil
}
func (s *rset) Del(keys ...any) error {
	return nil
}
func (s *rset) Items() []api.Value {
	return nil
}
