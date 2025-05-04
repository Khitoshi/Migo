package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Khitoshi/Migo.git/cmd"
	_ "github.com/lib/pq"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "migo",
	Short: "Migo is a database migration tool",
	Long:  `A simple and flexible database migration tool for PostgreSQL databases.`,
}

func init() {
	rootCmd.AddCommand(migrationCommand())
}

func migrationCommand() *cobra.Command {
	migrationCmd := &cobra.Command{
		Use:   "migration",
		Short: "Migration related commands",
		Long:  `Commands for creating and applying database migrations.`,
	}

	migrationCmd.AddCommand(cmd.MigrationInitCommand())
	migrationCmd.AddCommand(cmd.MigrationCreateCommand())
	migrationCmd.AddCommand(cmd.MigrationUpCommand())
	return migrationCmd
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
