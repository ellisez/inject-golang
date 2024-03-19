package model

type MultipleInstance struct {
	keys []string

	multipleMap map[string]*Provide
}

func newCtxMultipleInstance() *MultipleInstance {
	return &MultipleInstance{
		multipleMap: map[string]*Provide{},
	}
}

func (m *MultipleInstance) Add(data *Provide) {
	key := data.Instance
	m.Delete(key)
	m.keys = append(m.keys, key)
	m.multipleMap[key] = data
}

func (m *MultipleInstance) DeleteItem(data *Provide) {
	m.Delete(data.Instance)
}
func (m *MultipleInstance) Delete(key string) {
	if !m.Contains(key) {
		return
	}
	for i, k := range m.keys {
		if k == key {
			m.keys = append(m.keys[:i], m.keys[i+1:]...)
			break
		}
	}
	delete(m.multipleMap, key)
}

func (m *MultipleInstance) Contains(key string) bool {
	_, ok := m.multipleMap[key]
	return ok
}

func (m *MultipleInstance) Get(key string) *Provide {
	return m.multipleMap[key]
}

func (m *MultipleInstance) Replace(data *Provide) bool {
	key := data.Instance
	if !m.Contains(key) {
		m.Add(data)
		return false
	}
	m.multipleMap[key] = data
	return true
}
func (m *MultipleInstance) IndexOf(index int) *Provide {
	key := m.keys[index]
	return m.multipleMap[key]
}

// for sort

func (m *MultipleInstance) Len() int {
	return len(m.keys)
}

func (m *MultipleInstance) Swap(x int, y int) {
	old := m.keys[x]
	m.keys[x] = m.keys[y]
	m.keys[y] = old
}
func (m *MultipleInstance) Less(x int, y int) bool {
	a := m.IndexOf(x)
	b := m.IndexOf(y)
	orderA := a.Order
	orderB := b.Order
	if orderA != "" && orderB == "" {
		return true
	}
	return orderA < orderB
}
