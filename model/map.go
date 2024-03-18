package model

type OrderItem interface {
	GetOrder() string
	GetInstance() string
}
type OrderedMap[T OrderItem] struct {
	Keys []string
	Map  map[string]T
}

func NewOrderedMap[T OrderItem]() *OrderedMap[T] {
	return &OrderedMap[T]{Map: map[string]T{}}
}

// map operator

func (m *OrderedMap[T]) Add(data T) {
	key := data.GetInstance()
	m.Delete(key)
	m.Keys = append(m.Keys, key)
	m.Map[key] = data
}

func (m *OrderedMap[T]) DeleteItem(data T) {
	m.Delete(data.GetInstance())
}
func (m *OrderedMap[T]) Delete(key string) {
	if !m.Contains(key) {
		return
	}
	for i, k := range m.Keys {
		if k == key {
			m.Keys = append(m.Keys[:i], m.Keys[i+1:]...)
			break
		}
	}
	delete(m.Map, key)
}

func (m *OrderedMap[T]) Contains(key string) bool {
	_, ok := m.Map[key]
	return ok
}

func (m *OrderedMap[T]) Get(key string) T {
	return m.Map[key]
}

func (m *OrderedMap[T]) Replace(data T) bool {
	key := data.GetInstance()
	if !m.Contains(key) {
		m.Add(data)
		return false
	}
	m.Map[key] = data
	return true
}
func (m *OrderedMap[T]) IndexOf(index int) T {
	key := m.Keys[index]
	return m.Map[key]
}

// for sort

func (m *OrderedMap[T]) Len() int {
	return len(m.Keys)
}

func (m *OrderedMap[T]) Swap(x int, y int) {
	old := m.Keys[x]
	m.Keys[x] = m.Keys[y]
	m.Keys[y] = old
}
func (m *OrderedMap[T]) Less(x int, y int) bool {
	a := m.IndexOf(x)
	b := m.IndexOf(y)
	orderA := a.GetOrder()
	orderB := b.GetOrder()
	if orderA != "" && orderB == "" {
		return true
	}
	return orderA < orderB
}

// subtypes

type InstanceMap struct {
	*OrderedMap[Instance]
}

func NewInstanceMap() *InstanceMap {
	return &InstanceMap{OrderedMap: NewOrderedMap[Instance]()}
}

type ProxyMap struct {
	*OrderedMap[*Proxy]
}

func NewProxyMap() *ProxyMap {
	return &ProxyMap{OrderedMap: NewOrderedMap[*Proxy]()}
}
