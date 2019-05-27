package mysql

import (
	"database/sql"
	"database/sql/driver"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"
)

var dsn string

func init() {
	if os.Getenv("TEST_DSN") != "" {
		dsn = os.Getenv("TEST_DSN")
	} else {
		if !testing.Short() {
			panic("No test DSN defined")
		}
	}
}

func newTestDB(t *testing.T) (*sql.DB, func()) {

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		t.Fatal(err)
	}

	script, err := ioutil.ReadFile("./testdata/setup.sql")

	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	return db, func() {
		script, err := ioutil.ReadFile("./testdata/teardown.sql")
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
