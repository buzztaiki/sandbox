package util

func FromPtr[T any](p *T) T {
	if p == nil {
		var x T
		return x
	}
	return *p
}
