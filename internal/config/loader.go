package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Config holds key-value pairs parsed from an env-style config file.
type Config map[string]string

// LoadFromFile reads a .env style file and returns a Config map.
// Lines starting with '#' are treated as comments and ignored.
// Empty lines are skipped. Format: KEY=VALUE
func LoadFromFile(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer f.Close()

	cfg := make(Config)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: invalid format %q (expected KEY=VALUE)", lineNum, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" {
			return nil, fmt.Errorf("line %d: empty key", lineNum)
		}

		cfg[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return cfg, nil
}
