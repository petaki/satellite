package fake

import (
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/gomodule/redigo/redis"
)

// Redis is an in-memory stand-in for a Redis server, for use in tests.
type Redis struct {
	mu     sync.Mutex
	hashes map[string]map[string]string
	values map[string]string
	ttls   map[string]int
	calls  map[string]int
}

// NewRedis function.
func NewRedis() *Redis {
	return &Redis{
		hashes: map[string]map[string]string{},
		values: map[string]string{},
		ttls:   map[string]int{},
		calls:  map[string]int{},
	}
}

// HSet function.
func (f *Redis) HSet(key, field, value string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.hashes[key] == nil {
		f.hashes[key] = map[string]string{}
	}

	f.hashes[key][field] = value
}

// Get returns a plain string value and whether it was set.
func (f *Redis) Get(key string) (string, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()

	value, ok := f.values[key]

	return value, ok
}

// TTL returns the expiry set on a key, or zero when it has none.
func (f *Redis) TTL(key string) int {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.ttls[key]
}

// Keys returns every key the server holds, sorted.
func (f *Redis) Keys() []string {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.keys()
}

// Calls returns how many times a command was issued.
func (f *Redis) Calls(command string) int {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.calls[strings.ToUpper(command)]
}

// Pool function.
func (f *Redis) Pool() *redis.Pool {
	return &redis.Pool{
		MaxIdle: 1,
		Dial: func() (redis.Conn, error) {
			return &redisConn{redis: f}, nil
		},
	}
}

func (f *Redis) keys() []string {
	keys := make([]string, 0, len(f.hashes)+len(f.values))

	for key := range f.hashes {
		keys = append(keys, key)
	}

	for key := range f.values {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

func (f *Redis) do(command string, args []any) (any, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.calls[command]++

	switch command {
	case "PING":
		return "PONG", nil
	case "MULTI":
		return "OK", nil
	case "HGETALL":
		return f.hgetall(args)
	case "HMGET":
		return f.hmget(args)
	case "SCAN":
		return f.scan(args)
	case "EXISTS":
		return f.exists(args)
	case "SET":
		return f.set(args)
	case "EXPIRE":
		return f.expire(args)
	case "DEL":
		return f.del(args)
	}

	return nil, ErrUnknownCommand
}

func (f *Redis) del(args []any) (any, error) {
	key := arg(args, 0)

	_, hash := f.hashes[key]
	_, value := f.values[key]

	delete(f.hashes, key)
	delete(f.values, key)
	delete(f.ttls, key)

	if hash || value {
		return int64(1), nil
	}

	return int64(0), nil
}

func (f *Redis) exists(args []any) (any, error) {
	key := arg(args, 0)

	if _, ok := f.hashes[key]; ok {
		return int64(1), nil
	}

	if _, ok := f.values[key]; ok {
		return int64(1), nil
	}

	return int64(0), nil
}

func (f *Redis) expire(args []any) (any, error) {
	key := arg(args, 0)

	if _, ok := f.values[key]; !ok {
		if _, ok := f.hashes[key]; !ok {
			return int64(0), nil
		}
	}

	seconds, err := strconv.Atoi(arg(args, 1))
	if err != nil {
		return nil, err
	}

	f.ttls[key] = seconds

	return int64(1), nil
}

func (f *Redis) set(args []any) (any, error) {
	f.values[arg(args, 0)] = arg(args, 1)

	return "OK", nil
}

func (f *Redis) hgetall(args []any) (any, error) {
	hash := f.hashes[arg(args, 0)]
	fields := make([]string, 0, len(hash))

	for field := range hash {
		fields = append(fields, field)
	}

	sort.Strings(fields)

	reply := make([]any, 0, len(fields)*2)

	for _, field := range fields {
		reply = append(reply, []byte(field), []byte(hash[field]))
	}

	return reply, nil
}

func (f *Redis) hmget(args []any) (any, error) {
	hash := f.hashes[arg(args, 0)]
	reply := make([]any, 0, len(args)-1)

	for _, field := range args[1:] {
		value, ok := hash[toString(field)]
		if !ok {
			reply = append(reply, nil)

			continue
		}

		reply = append(reply, []byte(value))
	}

	return reply, nil
}

func (f *Redis) scan(args []any) (any, error) {
	cursor, _ := strconv.Atoi(arg(args, 0))
	match := "*"
	count := 10

	for i := 1; i+1 < len(args); i += 2 {
		switch strings.ToUpper(toString(args[i])) {
		case "MATCH":
			match = toString(args[i+1])
		case "COUNT":
			count, _ = strconv.Atoi(toString(args[i+1]))
		}
	}

	var matched []string

	for _, key := range f.keys() {
		if globMatch(match, key) {
			matched = append(matched, key)
		}
	}

	end := min(cursor+count, len(matched))

	next := end
	if next >= len(matched) {
		next = 0
	}

	page := make([]any, 0, end-cursor)

	for _, key := range matched[cursor:end] {
		page = append(page, []byte(key))
	}

	return []any{[]byte(strconv.Itoa(next)), page}, nil
}

type queued struct {
	command string
	args    []any
}

type redisConn struct {
	redis   *Redis
	pending []queued
}

func (c *redisConn) Close() error { return nil }

func (c *redisConn) Err() error { return nil }

func (c *redisConn) Do(command string, args ...any) (any, error) {
	if strings.ToUpper(command) == "EXEC" {
		return c.exec()
	}

	return c.redis.do(strings.ToUpper(command), args)
}

// Send queues a command, as redigo does inside a MULTI. An unsupported command
// fails here rather than silently doing nothing at EXEC time.
func (c *redisConn) Send(command string, args ...any) error {
	command = strings.ToUpper(command)

	if !supported(command) {
		return ErrUnknownCommand
	}

	c.pending = append(c.pending, queued{command: command, args: args})

	return nil
}

func (c *redisConn) exec() (any, error) {
	replies := make([]any, 0, len(c.pending))

	for _, p := range c.pending {
		reply, err := c.redis.do(p.command, p.args)
		if err != nil {
			c.pending = nil

			return nil, err
		}

		replies = append(replies, reply)
	}

	c.pending = nil
	c.redis.mu.Lock()
	c.redis.calls["EXEC"]++
	c.redis.mu.Unlock()

	return replies, nil
}

func supported(command string) bool {
	switch command {
	case "PING", "MULTI", "EXEC", "HGETALL", "HMGET", "SCAN", "EXISTS", "SET", "EXPIRE", "DEL":
		return true
	}

	return false
}

func (c *redisConn) Flush() error { return nil }

func (c *redisConn) Receive() (any, error) { return nil, nil }
