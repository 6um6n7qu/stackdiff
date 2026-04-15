package config

// MergeStrategy controls how conflicts are resolved when merging two configs.
type MergeStrategy int

const (
	// PreferBase keeps the base value when a key exists in both configs.
	PreferBase MergeStrategy = iota
	// PreferOverride replaces the base value with the override value.
	PreferOverride
)

// Merge combines base and override config maps into a new map.
// Keys present only in base or only in override are always included.
// Conflicting keys are resolved according to the provided MergeStrategy.
func Merge(base, override map[string]string, strategy MergeStrategy) map[string]string {
	result := make(map[string]string, len(base))

	for k, v := range base {
		result[k] = v
	}

	for k, v := range override {
		if _, exists := result[k]; !exists {
			result[k] = v
			continue
		}
		if strategy == PreferOverride {
			result[k] = v
		}
	}

	return result
}
