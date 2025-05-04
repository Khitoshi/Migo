package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func MigrationCreateCommand() *cobra.Command {
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
