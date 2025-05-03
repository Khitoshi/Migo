package main

import (
	"bufio"
	"os"
	"strings"
)

func LoadEnv(path string) (env map[string]string, err error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	env = make(map[string]string)

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}
		slice := strings.SplitN(line, "=", 2)
		if len(slice) != 2 {
			continue
		}

		key := strings.Trim(slice[0], " \"'`")
		rawValue := strings.TrimSpace(slice[1])

		// バックティックの数を数え、ちょうど2つの場合のみ左右を取り除く
		if strings.Count(rawValue, "`") == 2 && len(rawValue) >= 2 {
			first, last := rawValue[0], rawValue[len(rawValue)-1]
			if first == '`' && last == '`' {
				rawValue = rawValue[1 : len(rawValue)-1]
			}
		} else {
			// 外側のダブルクォート、シングルクォートは従来通り取り除く
			if len(rawValue) >= 2 {
				first, last := rawValue[0], rawValue[len(rawValue)-1]
				if (first == '"' && last == '"') ||
					(first == '\'' && last == '\'') {
					rawValue = rawValue[1 : len(rawValue)-1]
				}
			}
		}

		// 特殊ケース: シングルクォートの中の '' を単一の ' に置換
		if strings.HasPrefix(slice[1], "'") && strings.HasSuffix(slice[1], "'") {
			rawValue = strings.ReplaceAll(rawValue, "''", "'")
		}

		// 特殊ケース: バックティックの中の `` を単一の ` に置換
		if strings.HasPrefix(slice[1], "`") && strings.HasSuffix(slice[1], "`") {
			rawValue = strings.ReplaceAll(rawValue, "``", "`")
		}

		env[key] = rawValue
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return env, err
}
