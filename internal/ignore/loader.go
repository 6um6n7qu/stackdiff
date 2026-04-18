package ignore

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadFile reads an ignore pattern file (one pattern per line).
// Lines starting with '#' and blank lines are ignored.
func LoadFile(path string) (*List, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ignore: open %q: %w", path, err)
	}
	defer f.Close()

	var patterns []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ignore: scan %q: %w", path, err)
	}
	return New(patterns), nil
}
