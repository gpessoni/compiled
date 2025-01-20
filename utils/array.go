package utils

// Includes checks if an item is in a slice
func Includes[T comparable](slice []T, item T) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}

	return false
}

func Difference[T comparable](all, includes []T) []T {
	var diff []T

	for _, a := range all {
		if !Includes(includes, a) {
			diff = append(diff, a)
		}
	}

	return diff
}

// Remove Duplicate elements from a slice
func RemoveDuplicates[T comparable](elements []T) []T {
	encountered := map[T]bool{}
	result := []T{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}

	return result
}

// Find by function
func Find[T any](elements []T, f func(T) bool) (T, bool) {
	for _, element := range elements {
		if f(element) {
			return element, true
		}
	}

	var zero T
	return zero, false
}

// Filter by function
// If the function returns true, the element is included in the result
func Filter[T any](elements []T, f func(T) bool) []T {
	var filtered []T

	for _, element := range elements {
		if f(element) {
			filtered = append(filtered, element)
		}
	}

	return filtered
}
