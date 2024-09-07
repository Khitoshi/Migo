package database

import (
	"testing"
)

func TestDBConection(t *testing.T) {
	const dsn = "KHitoshi:985632@tcp(127.0.0.1:3306)/local_db"
	db, err := NewDB(dsn)
	if err != nil {
		t.Fatal(err)
	}

	if err = db.Close(); err != nil {
		t.Fatal(err)
	}
}
