package util

func FindElementInSlice(element string, target map[string]int) bool {
	if target[element] != 0 {
		return true
	}
	return false
}
