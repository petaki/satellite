package models

import (
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

const redisTestURLEnv = "REDIS_TEST_URL"

var commandStatPattern = regexp.MustCompile(`cmdstat_([a-z|]+):calls=(\d+)`)

type testRedis struct {
	t    *testing.T
	pool *redis.Pool
}

// newTestRedis connects to the server in REDIS_TEST_URL and empties its
// database, so point it at a throwaway instance only.
func newTestRedis(t *testing.T) *testRedis {
	t.Helper()

	url := os.Getenv(redisTestURLEnv)
	if url == "" {
		t.Skipf("%s is not set", redisTestURLEnv)
	}

	r := &testRedis{
		t: t,
		pool: &redis.Pool{
			MaxIdle: 3,
			Dial: func() (redis.Conn, error) {
				return redis.DialURL(url, redis.DialConnectTimeout(5*time.Second))
			},
		},
	}

	t.Cleanup(func() {
		r.pool.Close()
	})

	r.do("FLUSHDB")
	r.do("CONFIG", "RESETSTAT")

	return r
}

func (r *testRedis) Pool() *redis.Pool {
	return r.pool
}

func (r *testRedis) HSet(key, field, value string) {
	r.do("HSET", key, field, value)
}

func (r *testRedis) Get(key string) (string, bool) {
	value, err := redis.String(r.do("GET", key), nil)

	return value, err == nil
}

func (r *testRedis) TTL(key string) int {
	ttl, err := redis.Int(r.do("TTL", key), nil)
	if err != nil {
		r.t.Fatal(err)
	}

	return ttl
}

func (r *testRedis) Keys() []string {
	keys, err := redis.Strings(r.do("KEYS", "*"), nil)
	if err != nil {
		r.t.Fatal(err)
	}

	slices.Sort(keys)

	return keys
}

// Calls returns how many times the server ran a command since the test began.
func (r *testRedis) Calls(command string) int {
	info, err := redis.String(r.do("INFO", "commandstats"), nil)
	if err != nil {
		r.t.Fatal(err)
	}

	for _, match := range commandStatPattern.FindAllStringSubmatch(info, -1) {
		if match[1] == strings.ToLower(command) {
			calls, _ := strconv.Atoi(match[2])

			return calls
		}
	}

	return 0
}

func (r *testRedis) do(command string, args ...any) any {
	r.t.Helper()

	conn := r.pool.Get()
	defer conn.Close()

	reply, err := conn.Do(command, args...)
	if err != nil {
		r.t.Fatalf("%s: %v", command, err)
	}

	return reply
}
