// Package ptr holds tiny generic helpers for working with pointers.
//
// They exist because Go models "absent" as a pointer but gives no syntax for it: there is no
// optional type and no `?.`, so every crossing between an optional value and a nullable column is
// hand-written. What used to be the biggest source of that boilerplate — being unable to take the
// address of a value a function just returned — is gone as of Go 1.26, where the builtin `new`
// accepts an expression. Use `new(f(x))` directly; there is no `ptr.New` here.
package ptr

// NonZero returns a pointer to v, or nil when v is the zero value. It keeps "nobody" out of a
// nullable column: without it an absent caller reaches the database as the all-zeros uuid, which
// reads like a real user and cannot be told apart from one:
//
//	DeletedBy: ptr.NonZero(deletedBy),
func NonZero[T comparable](v T) *T {
	var zero T
	if v == zero {
		return nil
	}

	return &v
}

// Deref returns the pointed-to value, or the zero value when p is nil.
func Deref[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

// Map converts a pointer to another type while keeping nil as nil, so "not sent" survives the
// crossing instead of collapsing into a zero value:
//
//	TargetMinAge: ptr.Map(req.TargetAgeMin, toInt),
func Map[T, U any](p *T, convert func(T) U) *U {
	if p == nil {
		return nil
	}

	return new(convert(*p))
}
