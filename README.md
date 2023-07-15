# Redisson

[Redisson](https://github.com/redisson/redisson)-like wrapper 
for [mediocregopher/radix](https://github.com/mediocregopher/radix)

## Supported redis configurations
### Single-node
```go
client, err := redisson.NewConfig().
    WithName("TEST-SINGLE-CLIENT").
    WithDb(9).
    WithPoolSize(5).
    NewSingle(singleAddress)
```
### Cluster
```go
client, err := redisson.NewConfig().
    WithName("TEST-CLUSTER-CLIENT").
    WithPoolSize(5).
    NewCluster("127.0.0.1:7001", "127.0.0.1:7002", "127.0.0.1:7003")
```
### Sentinel
```go
client, err := redisson.NewConfig().
    WithName("TEST-SENTINEL-CLIENT").
    WithDb(9).
    WithPoolSize(5).
    NewSentinel("poolA", "127.0.0.1:26379", "127.0.0.1:26378", "127.0.0.1:26377")
```
### Authorized redis access
```go
client, err := redisson.NewConfig().
    WithName("TEST-SINGLE-CLIENT").
    WithAuth("login", "password").
    WithDb(9).
    WithPoolSize(5).
    NewSingle(singleAddress)
```
### Close
The connection should be closed after use
```go
err = client.Close()
```
## Data types
The data in redis keys is stored in "generic" form of Value interface.
```go
type Value interface {
	IsEmpty() bool
	String() string
	V() any
	AsString() string
	AsInt() int
	AsFloat() float64
	AsBool() bool
}
```
Basically it's translated to string representation and stored as string.

*![#f03c15](TODO)*: serializers for complex data structures still need to be implemented; 
for now only basic data types are supported:
- strings
- integer numbers
- real numbers
- booleans
## Supported redis functions
### Keyspace event notifications
Redis needs to be configured to send key-event notifications. 
This group of functions implements support for this functionality. 

Enable default ("KEAn") keyspace event notifications
```go
func EnableKeyEventNotifications() error
```

Enable keyspace event notifications of given types (see [redis documentation](https://redis.io/docs/manual/keyspace-notifications/#Configuration))
```go
func EnableKeyEventNotificationsOfTypes(types string) error
```

Disable keyspace event notifications
```go
func DisableKeyEventNotifications() error
```
### Common functions
	Del(keys ...string) (int, error)
	Expire(key string, ttl time.Duration) (int, error)
	Exists(key ...string) bool
	Keys(filter string) []string
	Touch(keys ...string)
	Type(key string) string
### Core
	Set(key string, value any) error    // Set set key value
	Get(key string) (Value, error)      // Get get key value
	Incr(key string) (int, error)       // Incr increment key value
	Decr(key string) (int, error)       // Decr decrement key value
### Collections
#### RList
	Len() int                           // Len returns list size
	LPush(items ...any) error           // LPush adds items to list tail in given order
	LPushRO(items ...any) error 	    // LPushRO adds items to list tail in reversed order
                                   	    //         i.e. first item in passed list will be added last
	LPop() (Value, error) 	            // LPop get item from list tail
	RPush(items ...any) error           // RPush adds items to list head in given order
	RPop() (Value, error)               // RPop get item from list head
#### RSet
	Size() int                          // Size return set size
	Add(value ...any) error             // Add adds items to the set
	Has(value any) bool                 // Has check is set has item
	Del(keys ...any) error              // Del removes items from the set
	Items() []Value                     // Items returns set items
#### RBitSet
	Set(idx uint32, value any) (bool, error)    // Set sets Nth bit of a set to passed value (0 / 1)
	Get(idx uint32) (bool, error)               // Get retrieves Nth bit of a set
	BitCount() int                              // BitCount returns number of set bits in bitset
	BitCountRange(start, end int, unit string) (int, error) // BitCountRange returns number of set bits 
                                                            //               in bitset on a given range
#### RMap
	Set(key string, value any) error    // set value for map key
	Get(key string) (Value, bool)       // get value of map key
	Del(keys ...string) error           // remove map element
	Keys() []string                     // retrieve a list of map keys
	Entries() []MapEntry                // retrieve a list of map entries
#### RCacheMap
Implements redis Map object with local cache. Runs background goroutine to synchronize local data to redis and back.

	Set(key string, value any) error    // set value for map key
	Get(key string) (Value, bool)       // get value of map key
	Del(keys ...string) error           // remove map element
	Keys() []string                     // retrieve a list of map keys
	Entries() []MapEntry                // retrieve a list of map entries
    Destroy()                           // destroy RCacheMap object
### PubSub
	PubSub() (radix.PubSubConn, error) // open pub-sub connection

Example use:
```go
psconn, err = client.PubSub()
if err != nil {
    return err
}
err = psconn.PSubscribe(context.Background(), fmt.Sprintf(keySpaceTopicFormat, m.key))
if err != nil {
    return err
}
ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
for {
    msg, err := m.psconn.Next(ctx)
    // process message
}
```

## TODO
- binary data codec support
- multi-node / multi-pool support


- simplified pub-sub
- streams


- bloomfilter
- cuckoofilter
- countminsketch
- json
- search
- tdigest
- timeseries
- topk


- lock
- rwlock
- distributed locks with redlock
