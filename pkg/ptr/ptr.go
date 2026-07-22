// Package ptr holds tiny generic helpers for working with pointers.
package ptr

// Deref returns the pointed-to value, or the zero value when p is nil.
func Deref[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}
