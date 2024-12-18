package database

import (
	"fmt"
)

/*
// migrationsテーブルが存在するか確認
func (db *DB) CheckMigrationTable() (bool, error) {
	const query = `SHOW TABLES WHERE Tables_in_mydb = ?;`
	var response string
	err := db.Conn.QueryRow(query).Scan(&response)
	if err != nil {
		return false, err
	}
	fmt.Printf("Check migration table:%s", response)
	return response != ``, nil
}

// migrationの一覧を取得
func (db *DB) GetMigrations() (string, error) {
	query := `SELECT table_name FROM migrations;`
	var response string
	err := db.Conn.QueryRow(query).Scan(&response)
	if err != nil {
		return "", err
	}
	fmt.Printf("Get migrations:%s", response)
	return response, nil
}
*/
// migrationsテーブルにレコードを追加
func (db *DB) AddMigration(table_name string) error {
	query := `INSERT INTO migrations (table_name) VALUES (?);`
	_, err := db.ExecuteQuery(query, table_name)
	if err != nil {
		return err
	}
	fmt.Println("Add migration")
	return nil
}
