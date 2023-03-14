package utils

func ToSliceOfAny[T any](s []T) []any {
	result := make([]any, len(s))
	for i, v := range s {
		result[i] = v
	}
	return result
}

func PageCount(total int, limit int) int {
	pages := total / limit

	if total%limit > 0 {
		pages++
	}

	return pages
}
