package utils

// ArrayContains checks if a string is in a given slice
func ArrayContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
