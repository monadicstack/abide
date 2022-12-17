package slices

// Remove is a generic way to remove the first instance of 'value' from the given slice.
// If the item does not exist in the slice, you'll get back your input slice as-is. This
// does not mutate your input slice - it returns a new slice without the value.
func Remove[T comparable](slice []T, value T) []T {
	for i, e := range slice {
		if e == value {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
