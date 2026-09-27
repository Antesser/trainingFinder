package utils

func TruncateForHasNext[T any](slice []T, limit int) ([]T, bool) {
	hasNext := len(slice) > limit
	if hasNext {
		slice = slice[:limit]
	}
	return slice, hasNext
}
