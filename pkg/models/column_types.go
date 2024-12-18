package models

import "fmt"

// ColumnType は SQL カラム型の基本定義
type ColumnType struct {
	SQLType       string
	Length        int
	Nullable      bool
	AutoIncrement bool
	PrimaryKey    bool
	Default       string
	Enum          []string
	IsBoolean     bool
}

// Column はテーブルのカラム定義
type Column struct {
	Name string
	Type ColumnType
}

// オプション型の定義
type ColumnOption func(*ColumnType)

// オプション関数の定義
func NotNull() ColumnOption {
	return func(col *ColumnType) {
		col.Nullable = false
	}
}

// オプション関数の定義
func AutoIncrement() ColumnOption {
	return func(col *ColumnType) {
		col.AutoIncrement = true
	}
}

func PrimaryKey() ColumnOption {
	return func(col *ColumnType) {
		col.PrimaryKey = true
	}
}

func Default(value string) ColumnOption {
	return func(col *ColumnType) {
		col.Default = value
	}
}

func Length(length int) ColumnOption {
	return func(col *ColumnType) {
		col.Length = length
	}
}

// カラム型の定義関数
func IntegerType(options ...ColumnOption) ColumnType {
	col := ColumnType{SQLType: "INT"}
	for _, option := range options {
		option(&col)
	}
	return col
}

func StringType(options ...ColumnOption) ColumnType {
	col := ColumnType{SQLType: "VARCHAR"}
	for _, option := range options {
		option(&col)
	}
	return col
}

func TimestampType(options ...ColumnOption) ColumnType {
	col := ColumnType{SQLType: "TIMESTAMP"}
	for _, option := range options {
		option(&col)
	}
	return col
}

func EnumType(values []string, options ...ColumnOption) ColumnType {
	col := ColumnType{SQLType: "ENUM", Enum: values}
	for _, option := range options {
		option(&col)
	}
	return col
}

func BooleanType(options ...ColumnOption) ColumnType {
	col := ColumnType{SQLType: "BOOLEAN"}
	for _, option := range options {
		option(&col)
	}
	return col
}

func CreateTableSQL(tableName string, columns []Column) string {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", tableName)

	for i, column := range columns {
		sql += fmt.Sprintf("%s %s", column.Name, column.Type.SQLType)

		// VARCHARの場合は長さを指定
		if column.Type.SQLType == "VARCHAR" && column.Type.Length > 0 {
			sql += fmt.Sprintf("(%d)", column.Type.Length)
		} else if column.Type.Length <= 0 {
			fmt.Printf("%s length is invalid", column.Name)
		}

		// ENUMの場合は値を指定
		if column.Type.SQLType == "ENUM" {
			sql += fmt.Sprintf("('%s')", join(column.Type.Enum, "', '"))
		}

		// Nullableオプション
		if !column.Type.Nullable {
			sql += " NOT NULL"
		} else {
			sql += " NULL"
		}

		// AutoIncrementオプション
		if column.Type.AutoIncrement {
			sql += " AUTO_INCREMENT"
		}

		// PrimaryKeyオプション
		if column.Type.PrimaryKey {
			sql += " PRIMARY KEY"
		}

		// Defaultオプション
		if column.Type.Default != "" {
			sql += fmt.Sprintf(" DEFAULT %s", column.Type.Default)
		}

		if i < len(columns)-1 {
			sql += ", "
		}

	}

	sql += ");"
	return sql
}

// ヘルパー関数
func join(elems []string, sep string) string {
	result := ""
	for i, elem := range elems {
		result += elem
		if i < len(elems)-1 {
			result += sep
		}
	}
	return result
}
