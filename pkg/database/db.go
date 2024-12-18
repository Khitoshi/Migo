package database

import (
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	Conn *sql.DB
}

// DBをOpen
func NewDB(dsn string) (*DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return &DB{Conn: db}, nil
}

// DBをクローズ
func (db *DB) Close() error {
	return db.Conn.Close()
}

// クエリを実行
func (db *DB) ExecuteQuery(query string, args ...interface{}) (interface{}, error) {
	query = strings.TrimSpace(strings.ToUpper(query))
	if strings.HasPrefix(query, "SELECT") {
		return db.executeSelectQuery(query, args...)
	}
	return db.executeNonSelectQuery(query, args...)
}

// SELECTクエリを実行
func (db *DB) executeSelectQuery(query string, args ...interface{}) (interface{}, error) {
	row := db.Conn.QueryRow(query, args...)
	var result interface{}
	err := row.Scan(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// SELECT以外のクエリを実行
func (db *DB) executeNonSelectQuery(query string, args ...interface{}) (interface{}, error) {
	result, err := db.Conn.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return result, nil
}
