package utils

// MapK return all keys as a slice
func MapK[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))

	for k := range m {
		r = append(r, k)
	}

	return r
}

// Source - https://stackoverflow.com/a/71635953
// Posted by blackgreen, modified by community. See post 'Timeline' for change history
// Retrieved 2026-10-12, License - CC BY-SA 4.0

// MapV return all values as a slice
func MapV[M ~map[K]V, K comparable, V any](m M) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}
