package utils

// Stores value of v and returns a pointer to it
func Bool(v bool) *bool {
	return &v
}

// Parses a string true as bool true
func ParseTrue(s string) bool {
	return s == "true"
}
