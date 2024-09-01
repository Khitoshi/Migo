package models

type ColumnType struct {
	SQLType string
}

var (
	SerialPrimaryKey = ColumnType{SQLType: "SERIAL PRIMARY KEY"}
	Varchar255       = ColumnType{SQLType: "VARCHAR(255)"}
	Varchar255Unique = ColumnType{SQLType: "VARCHAR(255) UNIQUE"}
)
