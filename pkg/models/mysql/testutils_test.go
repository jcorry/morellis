package mysql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func newTestDB(t *testing.T) (*sql.DB, func(), error) {
	var dsn string

	dsn = os.Getenv("TEST_DSN")
	if dsn == "" && !testing.Short() {
		t.Fatal("No test DSN defined")
	}

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		t.Fatal(err)
	}

	// Run DB migrations
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create DB driver for migrations: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", os.Getenv("MIGRATIONS_DIR")),
		"mysql",
		driver)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create migrator: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, nil, fmt.Errorf("unable to run DB migrations: %v", err)
	}

	// Do data inserts from testsdata dir
	script, err := ioutil.ReadFile("./testdata/setup.sql")

	if err != nil {
		t.Fatal(err)
	}

	queries := strings.Split(string(script), ";")
	tx, err := db.BeginTx(context.Background(), nil)

	if err != nil {
		return nil, nil, err
	}

	defer func() {
		_ = tx.Commit()
	}()

	for _, q := range queries {
		if q == "" {
			continue
		}
		_, err = tx.Exec(q)
		if err != nil {
			tx.Rollback()
			return nil, nil, err
		}
	}

	return db, func() {
		err = m.Drop()
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	}, nil
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
