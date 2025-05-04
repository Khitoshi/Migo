package cmd

import (
	"database/sql"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func MigrationInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize the database",
		Long:  `Initialize the database and create necessary tables for migration tracking.`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Initializing the database...")
			// 初期化ロジックをここに実装

			// root
			// └─migo
			//     └─migrations
			// directoryのチェック
			// もしなければ作成する
			// もしあれば、何もしない
			_, err := os.Stat("./migo")
			if os.IsNotExist(err) {
				fileInfo, err := os.Lstat("./")
				if err != nil {
					log.Println(err)
					return
				}

				fileMode := fileInfo.Mode()
				unixPerms := fileMode & os.ModePerm

				if err := os.MkdirAll("migo/migrations", unixPerms); err != nil {
					log.Fatal(err)
					return
				}

				log.Println("migo/migrations directory created.")
			} else {
				log.Println("migo/migrations directory already exists.")
			}

			// PostgreSQLへの接続
			db, err := sql.Open("postgres", "postgresql://khitoshi:985632@localhost:5432/postgres?sslmode=disable")
			if err != nil {
				panic(err)
			}
			defer db.Close()

			// migrationsテーブルの作成
			_, err = db.Exec(`
				CREATE TABLE IF NOT EXISTS migrations (
					id SERIAL PRIMARY KEY,
					name VARCHAR(255) NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);
			`)
			if err != nil {
				log.Fatal(err)
				return
			}
			log.Println("Migrations table created successfully.")

		},
	}
}
