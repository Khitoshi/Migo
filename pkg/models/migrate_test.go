package models

import (
	"testing"

	"github.com/Khitoshi/Migo/pkg/database"
)

func TestCreateTableSQL(t *testing.T) {
	const dsn = "KHitoshi:985632@tcp(127.0.0.1:3306)/local_db"
	db, err := database.NewDB(dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	tableName := "template"
	column := NewTemplateTable()
	sql := CreateTableSQL(tableName, column)
	t.Logf("CreateTableSQL:%s", sql)

	_, err = db.ExecuteQuery(sql)
	if err != nil {
		t.Fatal(err)
	}
}

func NewTemplateTable() []Column {
	return []Column{
		{"template_id", IntegerType(NotNull(), AutoIncrement(), PrimaryKey())},
		{"template_title", StringType(Length(64), NotNull())},
		{"template_text", StringType(Length(2000), NotNull())},
		{"language", EnumType([]string{"jp", "en", "cn"}, NotNull())},
		{"is_default", BooleanType(NotNull(), Default("0"))},
		{"created_at", TimestampType(NotNull(), Default("CURRENT_TIMESTAMP"))},
		{"updated_at", TimestampType(NotNull(), Default("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"))},
	}
}
