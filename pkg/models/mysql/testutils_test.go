package mysql

import (
	"database/sql"
	"flag"
	"io/ioutil"
	"testing"
)

var dsn string

func init() {
	flag.StringVar(&dsn, "dsn", "morellistest:testpass@tcp(127.0.0.1:33062)/morellistest?parseTime=true&multiStatements=true", "MySQL DSN URL")
	flag.Parse()
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
