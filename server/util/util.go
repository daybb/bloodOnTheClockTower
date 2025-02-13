package util

func FindElementInSlice(element int, slice []int) int {
	for i := range slice {
		if slice[i] == element {
			return i
		}
	}
	return -1
}
