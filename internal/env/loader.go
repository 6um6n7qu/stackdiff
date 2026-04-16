package env

import (
	"bufio"
	"os"
	"strings"
)

// LoadFile reads a .env file and returns a map of key-value pairs.
// Lines starting with '#' and empty lines are ignored.
// Values may optionally be quoted with double quotes.
func LoadFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := make(map[string]string)
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := parseLine(line)
		if !ok {
			continue
		}
		result[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

// parseLine splits a line into key and value on the first '=' character.
// Strips surrounding double quotes from the value if present.
func parseLine(line string) (string, string, bool) {
	idx := strings.IndexByte(line, '=')
	if idx < 1 {
		return "", "", false
	}
	key := strings.TrimSpace(line[:idx])
	val := strings.TrimSpace(line[idx+1:])
	if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
		val = val[1 : len(val)-1]
	}
	return key, val, true
}
