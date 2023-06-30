package redisson

type RSet interface {
	Add(value any) error
	Has(value any) (bool, error)
	Del(keys ...any) error
	Items() []Value
}

type rset struct {
	client Redis
	key    string
}

func NewRSet(key string, client Redis) RSet {
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
func (s *rset) Items() []Value {
	return nil
}
