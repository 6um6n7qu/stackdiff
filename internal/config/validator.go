package config

import "fmt"

// ValidationError holds a list of issues found during validation.
type ValidationError struct {
	Issues []string
}

func (e *ValidationError) Error() string {
	msg := fmt.Sprintf("%d validation issue(s) found:", len(e.Issues))
	for _, issue := range e.Issues {
		msg += "\n  - " + issue
	}
	return msg
}

// Validate checks that a config map is non-nil, has at least one key,
// and that no key or value is an empty string.
func Validate(cfg map[string]string) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	var issues []string

	if len(cfg) == 0 {
		issues = append(issues, "config contains no keys")
	}

	for k, v := range cfg {
		if k == "" {
			issues = append(issues, "empty key found")
			continue
		}
		if v == "" {
			issues = append(issues, fmt.Sprintf("key %q has an empty value", k))
		}
	}

	if len(issues) > 0 {
		return &ValidationError{Issues: issues}
	}
	return nil
}

// MustValidate is like Validate but panics if the config is invalid.
// It is intended for use during program initialization where an invalid
// config is an unrecoverable error.
func MustValidate(cfg map[string]string) {
	if err := Validate(cfg); err != nil {
		panic("config validation failed: " + err.Error())
	}
}
