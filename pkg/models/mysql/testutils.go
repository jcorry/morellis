package mysql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func NewTestDB(t *testing.T) (*sql.DB, func()) {
	var dsn string

	dsn = os.Getenv("TEST_DSN")
	if dsn == "" && !testing.Short() {
		t.Fatal("No test DSN defined")
	}

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		t.Fatal(err)
	}

	script, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", os.Getenv("TEST_DATA_DIR"), "setup.sql"))

	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	return db, func() {
		script, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", os.Getenv("TEST_DATA_DIR"), "teardown.sql"))
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	}
}

func NewTestRedis(t *testing.T) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("TEST_REDIS_ADDRESS"),
		Password: "",
		DB:       0,
	})

	cmd := rdb.Ping(context.TODO())
	if cmd.Err() != nil {
		t.Fatal(cmd.Err())
	}
	t.Log(fmt.Sprintf("Redis: %s", cmd.String()))

	t.Cleanup(func() {
		rdb.FlushDB(context.TODO())
	})

	return rdb
}

func randomTimestamp() time.Time {
	randomTime := rand.Int63n(time.Now().Unix()-94608000) + 94608000

	randomNow := time.Unix(randomTime, 0)

	return randomNow
}

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type AnyInt64 struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyInt64) Match(v driver.Value) bool {
	_, ok := v.(int64)
	return ok
}
