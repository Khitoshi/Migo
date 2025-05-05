package cmd

import (
	"database/sql"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	_ "github.com/lib/pq" // PostgreSQLドライバ
	"github.com/spf13/cobra"
)

func MigrationUpCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Apply all pending migrations",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Applying all pending migrations...")

			// PostgreSQLへの接続
			db, err := sql.Open("postgres", "postgresql://khitoshi:985632@localhost:5432/postgres?sslmode=disable")
			if err != nil {
				panic(err)
			}
			defer db.Close()

			// マイグレーションテーブルの存在確認、なければ作成
			_, err = db.Exec(`
                CREATE TABLE IF NOT EXISTS migrations (
                    id SERIAL PRIMARY KEY,
                    name VARCHAR(255) NOT NULL,
                    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
                )
            `)
			if err != nil {
				log.Fatalf("Failed to create migrations table: %v", err)
			}

			// migrationsディレクトリのファイルを読み込む
			migrationsDir := "./migo/migrations"
			files, err := os.ReadDir(migrationsDir)
			if err != nil {
				log.Fatalf("Error reading directory: %v", err)
			}

			for _, f := range files {
				if !f.IsDir() && strings.HasSuffix(f.Name(), ".go") {
					filePath := filepath.Join(migrationsDir, f.Name())
					log.Printf("Processing file: %s", filePath)

					// ファイルの情報を取得
					fileInfo, err := f.Info()
					if err != nil {
						log.Fatalf("Error getting file info: %v", err)
					}

					// マイグレーション情報を確認
					var updatedAt time.Time
					rows, err := db.Query("SELECT updated_at FROM migrations WHERE name = $1", f.Name())
					if err != nil {
						log.Printf("Error checking migration status: %v", err)
						continue
					}

					if rows.Next() {
						err = rows.Scan(&updatedAt)
						rows.Close()
						if err != nil {
							log.Printf("Error scanning updated_at: %v", err)
							continue
						}
					} else {
						rows.Close()
						// 新しいマイグレーション: 登録
						_, err = db.Exec("INSERT INTO migrations (name, updated_at) VALUES ($1, $2)",
							f.Name(), fileInfo.ModTime())
						if err != nil {
							log.Printf("Error inserting migration record: %v", err)
							continue
						}
					}

					// 変更がない場合はスキップ
					if !updatedAt.IsZero() && !fileInfo.ModTime().After(updatedAt) {
						log.Printf("No changes for %s, skipping", f.Name())
						continue
					}

					// ファイルから構造体を解析
					structDefs, err := parseStructsFromFile(filePath)
					if err != nil {
						log.Printf("Error parsing file %s: %v", f.Name(), err)
						continue
					}

					// 各構造体に対してテーブル作成
					for _, structDef := range structDefs {
						log.Printf("Creating table for struct: %s", structDef.Name)

						// テーブル名は構造体名を小文字に
						tableName := strings.ToLower(structDef.Name)

						// テーブル作成SQLを生成
						createTableSQL := generateCreateTableSQL(tableName, structDef)
						log.Printf("Generated SQL: %s", createTableSQL)

						// SQLを実行
						_, err = db.Exec(createTableSQL)
						if err != nil {
							log.Printf("Error creating table %s: %v", tableName, err)
							continue
						}

						log.Printf("Successfully created table: %s", tableName)
					}

					// マイグレーション情報を更新
					if !updatedAt.IsZero() {
						_, err = db.Exec("UPDATE migrations SET updated_at = $1 WHERE name = $2",
							fileInfo.ModTime(), f.Name())
						if err != nil {
							log.Printf("Error updating migration record: %v", err)
						}
					}
				}
			}
		},
	}
}

// StructField は構造体のフィールド情報を保持します
type StructField struct {
	Name    string
	Type    string
	Tags    map[string]string
	Options []string
}

// StructDef は構造体定義を表します
type StructDef struct {
	Name   string
	Fields []StructField
}

// parseStructsFromFile はGoソースファイルから構造体定義を抽出します
func parseStructsFromFile(filePath string) ([]StructDef, error) {
	// ファイルセットを作成
	fset := token.NewFileSet()

	// ファイルを解析
	f, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("ファイルの解析に失敗: %w", err)
	}

	var structDefs []StructDef

	// ファイル内の宣言をループ
	for _, decl := range f.Decls {
		// 一般的な宣言のみを処理
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		// 型宣言のみを処理
		if genDecl.Tok != token.TYPE {
			continue
		}

		// 宣言内の仕様をループ
		for _, spec := range genDecl.Specs {
			// 型仕様のみを処理
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			// 構造体のみを処理
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			// 構造体定義を作成
			structDef := StructDef{
				Name: typeSpec.Name.Name,
			}

			// フィールドをループ
			for _, field := range structType.Fields.List {
				if len(field.Names) == 0 {
					continue // 無名フィールドをスキップ
				}

				// フィールド情報を取得
				fieldName := field.Names[0].Name
				fieldType := ""

				// フィールドの型を文字列として取得
				switch t := field.Type.(type) {
				case *ast.Ident:
					fieldType = t.Name
				case *ast.SelectorExpr:
					fieldType = fmt.Sprintf("%s.%s", t.X, t.Sel.Name)
				default:
					fieldType = fmt.Sprintf("%T", field.Type)
				}

				// タグを解析
				tags := make(map[string]string)
				options := []string{}

				if field.Tag != nil {
					tag := field.Tag.Value

					// タグからクォートを削除
					tag = strings.Trim(tag, "`")

					// タグを空白で分割
					structTag := reflect.StructTag(tag)

					// json タグを取得
					jsonTag := structTag.Get("json")
					if jsonTag != "" {
						tags["json"] = jsonTag
					}

					// opt タグを取得
					optTag := structTag.Get("opt")
					if optTag != "" {
						options = strings.Split(optTag, ",")
						tags["opt"] = optTag
					}
				}

				// フィールド情報を構造体定義に追加
				structDef.Fields = append(structDef.Fields, StructField{
					Name:    fieldName,
					Type:    fieldType,
					Tags:    tags,
					Options: options,
				})
			}

			// 構造体定義をリストに追加
			structDefs = append(structDefs, structDef)
		}
	}

	return structDefs, nil
}

// generateCreateTableSQL は構造体定義からCREATE TABLE SQLを生成します
func generateCreateTableSQL(tableName string, structDef StructDef) string {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", tableName)

	for i, field := range structDef.Fields {
		columnName := field.Name
		columnType := "TEXT" // デフォルト型
		columnOptions := ""

		// JSON タグからカラム名を取得
		if jsonTag, exists := field.Tags["json"]; exists {
			parts := strings.Split(jsonTag, ":")
			if len(parts) > 1 {
				columnName = strings.Trim(parts[1], `"`)
			}
		}

		// フィールドの型をSQLタイプに変換
		switch field.Type {
		case "int", "int32", "int64":
			columnType = "INTEGER"
		case "string":
			columnType = "VARCHAR(255)"
		case "float32", "float64":
			columnType = "NUMERIC"
		case "bool":
			columnType = "BOOLEAN"
		case "time.Time":
			columnType = "TIMESTAMP"
		}

		// オプションを処理
		hasPrimaryKey := false
		hasAutoIncrement := false

		for _, opt := range field.Options {
			switch opt {
			case "primary_key":
				columnOptions += " PRIMARY KEY"
				hasPrimaryKey = true
			case "auto_increment":
				hasAutoIncrement = true
			case "not_null":
				columnOptions += " NOT NULL"
			}
		}

		// PostgreSQL の SERIAL 型を処理
		if hasAutoIncrement && columnType == "INTEGER" && hasPrimaryKey {
			columnType = "SERIAL"
			// PRIMARY KEYはSERIAL型に含まれるので削除
			columnOptions = strings.Replace(columnOptions, " PRIMARY KEY", "", 1)
			columnOptions += " PRIMARY KEY"
		}

		// カラム定義を追加
		sql += fmt.Sprintf("    %s %s%s", columnName, columnType, columnOptions)

		// カンマを追加（最後のフィールド以外）
		if i < len(structDef.Fields)-1 {
			sql += ",\n"
		}
	}

	sql += "\n);"
	return sql
}

func insertData(db *sql.DB, name string) error {
	result, err := db.Exec("INSERT INTO migrations (name) VALUES ($1)", name)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	log.Printf("Inserted record ID: %d\n", id)
	return nil
}
