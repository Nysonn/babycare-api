package parent

import "strconv"

// parseFloat converts a string representation of a numeric DB value to float64.
// Returns 0 if the string is empty or cannot be parsed.
func parseFloat(s string) float64 {
	if s == "" {
		return 0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}
