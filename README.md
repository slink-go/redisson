# Redisson

Basic [Redisson](https://github.com/redisson/redisson)-like wrapper 
for [mediocregopher/radix](https://github.com/mediocregopher/radix)

# Table of Contents
1. [Redis Configurations](#connection)
   - [Single node](#single.node.connection)
   - [Cluster](#clustered.connection)
   - [Sentinel](#sentinel.connection)
   - [Authorization](#authorized.access.connection)
   - [Close](#close.connection)
2. [Data types](#data-types)
3. [Supported redis functions](#supported.functions)
   - [Keyspace event notifications](#supported.functions.ksn)
   - [Common functions](#supported.functions.common)
   - [Core functions](#supported.functions.core)
   - [Collections](#supported.functions.collections)
     - [RList](#supported.functions.collections.rlist)
     - [RSet](#supported.functions.collections.rset)
     - [RBitSet](#supported.functions.collections.rbitset)
     - [RMap](#supported.functions.collections.rmap)
     - [RCacheMap](#supported.functions.collections.rcachemap)
4. [PubSub](#supported.functions.pubsub)
5. [TODO](#todo)


## Supported redis configurations<a name="connection"></a>
### Single-node<a name="single.node.connection"></a>
```go
client, err := redisson.NewConfig().
    WithName("TEST-SINGLE-CLIENT").
    WithDb(9).
    WithPoolSize(5).
    NewSingle(singleAddress)
```
### Cluster<a name="clustered.connection"></a>
```go
client, err := redisson.NewConfig().
    WithName("TEST-CLUSTER-CLIENT").
    WithPoolSize(5).
    NewCluster("127.0.0.1:7001", "127.0.0.1:7002", "127.0.0.1:7003")
```
### Sentinel<a name="sentinel.connection"></a>
```go
client, err := redisson.NewConfig().
    WithName("TEST-SENTINEL-CLIENT").
    WithDb(9).
    WithPoolSize(5).
    NewSentinel("poolA", "127.0.0.1:26379", "127.0.0.1:26378", "127.0.0.1:26377")
```
### Authorized redis access<a name="authorized.access.connection"></a>
```go
client, err := redisson.NewConfig().
    WithName("TEST-SINGLE-CLIENT").
    WithAuth("login", "password").
    WithDb(9).
    WithPoolSize(5).
    NewSingle(singleAddress)
```
### Close<a name="close.connection"></a>
The connection should be closed after use
```go
err = client.Close()
```
## Data types<a name="data-types"></a>
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

> __Warning__
<br>TODO: serializers for complex data structures still need to be implemented; 
for now only basic data types are supported:
<br>- strings
<br>- integer numbers
<br>- real numbers
<br>- booleans
## Supported redis functions<a name="supported.functions"></a>
### Keyspace event notifications<a name="supported.functions.ksn"></a>
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
### Common functions<a name="supported.functions.common"></a>
	Del(keys ...string) (int, error)
	Expire(key string, ttl time.Duration) (int, error)
	Exists(key ...string) bool
	Keys(filter string) []string
	Touch(keys ...string)
	Type(key string) string
### Core<a name="supported.functions.core"></a>
	Set(key string, value any) error    // Set set key value
	Get(key string) (Value, error)      // Get get key value
	Incr(key string) (int, error)       // Incr increment key value
	Decr(key string) (int, error)       // Decr decrement key value
### Collections<a name="supported.functions.collections"></a>
#### RList<a name="supported.functions.collections.rlist"></a>
	Len() int                           // Len returns list size
	LPush(items ...any) error           // LPush adds items to list tail in given order
	LPushRO(items ...any) error 	    // LPushRO adds items to list tail in reversed order
                                   	    //         i.e. first item in passed list will be added last
	LPop() (Value, error) 	            // LPop get item from list tail
	RPush(items ...any) error           // RPush adds items to list head in given order
	RPop() (Value, error)               // RPop get item from list head
#### RSet<a name="supported.functions.collections.rset"></a>
	Size() int                          // Size return set size
	Add(value ...any) error             // Add adds items to the set
	Has(value any) bool                 // Has check is set has item
	Del(keys ...any) error              // Del removes items from the set
	Items() []Value                     // Items returns set items
#### RBitSet<a name="supported.functions.collections.rbitset"></a>
	Set(idx uint32, value any) (bool, error)    // Set sets Nth bit of a set to passed value (0 / 1)
	Get(idx uint32) (bool, error)               // Get retrieves Nth bit of a set
	BitCount() int                              // BitCount returns number of set bits in bitset
	BitCountRange(start, end int, unit string) (int, error) // BitCountRange returns number of set bits 
                                                            //               in bitset on a given range
#### RMap<a name="supported.functions.collections.rmap"></a>
	Set(key string, value any) error    // set value for map key
	Get(key string) (Value, bool)       // get value of map key
	Del(keys ...string) error           // remove map element
	Keys() []string                     // retrieve a list of map keys
	Entries() []MapEntry                // retrieve a list of map entries
#### RCacheMap<a name="supported.functions.collections.rcachemap"></a>
Implements redis Map object with local cache. Runs background goroutine to synchronize local data to redis and back.

	Set(key string, value any) error    // set value for map key
	Get(key string) (Value, bool)       // get value of map key
	Del(keys ...string) error           // remove map element
	Keys() []string                     // retrieve a list of map keys
	Entries() []MapEntry                // retrieve a list of map entries
    Destroy()                           // destroy RCacheMap object
### PubSub<a name="supported.functions.pubsub"></a>
	PubSub() (radix.PubSubConn, error) // open pub-sub connection

Example use:
```go
psconn, err = client.PubSub()
if err != nil {
    return err
}
err = psconn.PSubscribe(context.Background(), "subscription-topic")
if err != nil {
    return err
}
ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
for {
    msg, err := m.psconn.Next(ctx)
    // process message
}
```
## TODO<a name="todo"></a>
```
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
```
