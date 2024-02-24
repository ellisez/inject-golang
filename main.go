package inject

var Container = make(map[string]any)

func Get[T any](name string) T {
	value := Container[name]
	return value.(T)
}
func Set[T any](name string, value T) {
	Container[name] = value
}
