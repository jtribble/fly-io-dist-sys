package slices

func Contains[V comparable](values []V, value V) bool {
	for _, s := range values {
		if s == value {
			return true
		}
	}
	return false
}
