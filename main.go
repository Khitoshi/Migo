package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "migo",
	Short: "Migo is a database migration tool",
	Long:  `A simple and flexible database migration tool for PostgreSQL databases.`,
}

func init() {
	rootCmd.AddCommand(initCommand())
	rootCmd.AddCommand(migrationCommand())
}

func initCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize the database",
		Long:  `Initialize the database and create necessary tables for migration tracking.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Initializing the database...")
			// 初期化ロジックをここに実装
		},
	}
}

func migrationCommand() *cobra.Command {
	migrationCmd := &cobra.Command{
		Use:   "migration",
		Short: "Migration related commands",
		Long:  `Commands for creating and applying database migrations.`,
	}

	migrationCmd.AddCommand(createMigrationCommand())
	migrationCmd.AddCommand(upMigrationCommand())

	return migrationCmd
}

func createMigrationCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new migration file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			fmt.Printf("Creating migration file: %s\n", name)
			// マイグレーションファイル作成ロジックをここに実装
		},
	}
}

func upMigrationCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Apply all pending migrations",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Applying all pending migrations...")
			// マイグレーション適用ロジックをここに実装
		},
	}
}

func main() {
	// .env読み込み
	env, err := LoadEnv(".env")
	if err != nil {
		panic(err)
	}
	for key, value := range env {
		println(key, "=", value)
	}

	// PostgreSQLへの接続
	db, err := sql.Open("postgres", env["DATABASE_URL"])
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 接続テスト
	if err := db.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "Could not ping database: %v\n", err)
		os.Exit(1)
	}

	// コマンド実行
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
		os.Exit(1)
	}

}
