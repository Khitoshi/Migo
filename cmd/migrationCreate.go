package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/spf13/cobra"
)

func toCamelCase(s string) string {
	// ハイフンやアンダースコアで分割
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '-' || r == '_' || unicode.IsSpace(r)
	})

	if len(words) == 0 {
		return ""
	}

	var result string
	// 以降の単語は先頭を大文字に
	for _, word := range words {
		if len(word) > 0 {
			result += strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}

	return result
}

func MigrationCreateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new migration file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			fmt.Printf("Creating migration file: %s\n", name)

			if len(name) == 0 {
				log.Fatal("Migration name cannot be empty")
			}
			if strings.Contains(name, " ") {
				log.Fatal("Migration name cannot contain spaces")
			}

			// マイグレーションファイル作成ロジックをここに実装
			f, err := os.Create("./migo/migrations/" + name + ".go")
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			_, err = f.WriteString(
				"// Migration: " + toCamelCase(name) + "\n" +
					"package migrations\n" +
					"\n" +
					"type " + toCamelCase(name) + " struct {\n}")
			if err != nil {
				log.Fatal(err)
			}
		},
	}
}
