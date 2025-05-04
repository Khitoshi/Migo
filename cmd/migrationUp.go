package cmd

import (
	"database/sql"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func MigrationUpCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Apply all pending migrations",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Applying all pending migrations...")
			// マイグレーション適用ロジックをここに実装

			files, err := os.ReadDir("./migo/migrations")
			if err != nil {
				log.Fatalf("Error reading directory: %v", err)
			}
			for _, f := range files {
				if !f.IsDir() {
					//ファイル名と更新日時を表示
					log.Println("File Name:", f.Name())
					//ファイルの更新日時を取得
					fileInfo, err := f.Info()
					if err != nil {
						log.Fatalf("Error getting file info: %v", err)
					}
					//更新日時を表示
					log.Println("Last Modified:", fileInfo.ModTime())

					//migrationテーブルにnameが登録されているかを確認
					//もし登録されていなければ、migrationテーブルにnameを登録する
					//もし登録されていれば、更新日時を比較して、更新日時が新しければ、migrationテーブルのupdated_atを更新とDBにINSERTする

					// PostgreSQLへの接続
					db, err := sql.Open("postgres", "postgresql://khitoshi:985632@localhost:5432/postgres?sslmode=disable")
					if err != nil {
						panic(err)
					}
					defer db.Close()

					//structを取得して、メンバ変数に設定されているjsonからカラム名とoptionを取得
					//それを使用してテーブルとカラムを作成する

				}
			}
		},
	}
}
