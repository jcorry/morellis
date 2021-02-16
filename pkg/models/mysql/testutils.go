package mysql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewTestDB(t *testing.T) *sql.DB {
	var dsn string

	dsn = os.Getenv("TEST_DSN")
	if dsn == "" && !testing.Short() {
		t.Fatal("No test DSN defined")
	}

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		t.Fatal(err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		t.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", os.Getenv("MIGRATIONS_DIR")),
		"mysql",
		driver,
	)
	if err != nil {
		t.Fatal(err)
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err = m.Drop()
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	})

	return db
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
