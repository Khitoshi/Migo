package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func MigrationUpCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Apply all pending migrations",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Applying all pending migrations...")
			// マイグレーション適用ロジックをここに実装
		},
	}
}
