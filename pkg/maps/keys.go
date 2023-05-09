package maps

func Keys[K comparable, V any](values map[K]V) []K {
	keys := make([]K, len(values))
	i := 0
	for k := range values {
		keys[i] = k
		i += 1
	}
	return keys
}
