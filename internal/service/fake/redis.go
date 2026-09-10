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
	calls  map[string]int
}

// NewRedis function.
func NewRedis() *Redis {
	return &Redis{
		hashes: map[string]map[string]string{},
		values: map[string]string{},
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

// Set function.
func (f *Redis) Set(key, value string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.values[key] = value
}

// Keys function.
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
	case "HGETALL":
		return f.hgetall(args)
	case "HMGET":
		return f.hmget(args)
	case "HSET":
		return f.hset(args)
	case "EXISTS":
		return f.exists(args)
	case "DEL":
		return f.del(args)
	case "SET":
		f.values[arg(args, 0)] = arg(args, 1)

		return "OK", nil
	case "EXPIRE", "MULTI", "EXEC":
		return "OK", nil
	case "SCAN":
		return f.scan(args)
	}

	return nil, ErrUnknownCommand
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

func (f *Redis) hset(args []any) (any, error) {
	key := arg(args, 0)

	if f.hashes[key] == nil {
		f.hashes[key] = map[string]string{}
	}

	for i := 1; i+1 < len(args); i += 2 {
		f.hashes[key][toString(args[i])] = toString(args[i+1])
	}

	return int64(1), nil
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

func (f *Redis) del(args []any) (any, error) {
	for _, key := range args {
		delete(f.hashes, toString(key))
		delete(f.values, toString(key))
	}

	return int64(len(args)), nil
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

type redisConn struct {
	redis  *Redis
	queued [][]any
}

func (c *redisConn) Close() error { return nil }

func (c *redisConn) Err() error { return nil }

func (c *redisConn) Do(command string, args ...any) (any, error) {
	if strings.ToUpper(command) == "EXEC" {
		for _, queued := range c.queued {
			if _, err := c.redis.do(toString(queued[0]), queued[1:]); err != nil {
				return nil, err
			}
		}

		c.queued = nil

		return []any{}, nil
	}

	return c.redis.do(strings.ToUpper(command), args)
}

func (c *redisConn) Send(command string, args ...any) error {
	if strings.ToUpper(command) == "MULTI" {
		c.queued = nil

		return nil
	}

	c.queued = append(c.queued, append([]any{command}, args...))

	return nil
}

func (c *redisConn) Flush() error { return nil }

func (c *redisConn) Receive() (any, error) { return nil, nil }
