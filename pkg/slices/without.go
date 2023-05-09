package slices

func Without[V comparable](values []V, value V) []V {
	newValues := make([]V, len(values))
	i := 0
	for _, s := range values {
		if s != value {
			newValues[i] = s
			i += 1
		}
	}
	return newValues[:i]
}
