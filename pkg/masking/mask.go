package masking

import "strings"

// MaskString masks all but the last five characters of a string with '*'.
func MaskString(s string) string {
	// Check if the string length is less than or equal to 4
	if len(s) <= 5 {
		return s
	}

	// Replace all but the last four characters with '*'
	masked := strings.Repeat("*", len(s)-5) + s[len(s)-5:]

	return masked
}
