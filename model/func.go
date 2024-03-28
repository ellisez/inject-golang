package model

type FuncInstance struct {
	keys []*Key

	funcMap map[string]*Proxy

	methodMap map[string]*Proxy
}

func newCtxFuncInstance() *FuncInstance {
	return &FuncInstance{
		funcMap:   map[string]*Proxy{},
		methodMap: map[string]*Proxy{},
	}
}
func (f *FuncInstance) AddFunc(data *Proxy) {
	key := data.Instance
	f.Delete(key)
	f.keys = append(f.keys, &Key{
		Instance: key,
		Type:     "func",
		Order:    data.Order,
	})
	f.funcMap[key] = data
}
func (f *FuncInstance) AddMethod(data *Proxy) {
	key := data.Instance
	f.Delete(key)
	f.keys = append(f.keys, &Key{
		Instance: key,
		Type:     "method",
		Order:    data.Order,
	})
	f.methodMap[key] = data
}
func (f *FuncInstance) Add(data *Proxy) {
	if data.Recv == nil {
		f.AddFunc(data)
	} else {
		f.AddMethod(data)
	}
}

func (f *FuncInstance) Delete(key string) {
	if !f.Contains(key) {
		return
	}
	var hitKey *Key
	for i, k := range f.keys {
		if k.Instance == key {
			f.keys = append(f.keys[:i], f.keys[i+1:]...)
			hitKey = k
			break
		}
	}
	switch hitKey.Type {
	case "func":
		delete(f.funcMap, key)
	case "method":
		delete(f.methodMap, key)
	}

}

func (f *FuncInstance) Contains(key string) bool {
	return f.Get(key) != nil
}
func (f *FuncInstance) GetFunc(key string) *Proxy {
	return f.funcMap[key]
}
func (f *FuncInstance) GetMethod(key string) *Proxy {
	return f.methodMap[key]
}
func (f *FuncInstance) Get(key string) *Proxy {
	fn := f.funcMap[key]
	if fn != nil {
		return fn
	}
	method := f.methodMap[key]
	if method != nil {
		return method
	}
	return nil
}

func (f *FuncInstance) Replace(data *Proxy) bool {
	key := data.Instance
	k := f.KeyOf(key)
	if k == nil {
		return false
	}
	switch k.Type {
	case "func":
		delete(f.funcMap, key)
	case "method":
		delete(f.methodMap, key)
	}
	k.Order = data.Order
	if data.Recv == nil {
		k.Type = "func"
		f.funcMap[key] = data
	} else {
		k.Type = "method"
		f.methodMap[key] = data
	}
	return true
}
func (f *FuncInstance) IndexOf(index int) *Proxy {
	key := f.keys[index]
	return f.Get(key.Instance)
}
func (f *FuncInstance) KeyOf(key string) *Key {
	for _, k := range f.keys {
		if k.Instance == key {
			return k
		}
	}
	return nil
}

func (f *FuncInstance) FuncLen() int {
	return len(f.funcMap)
}
func (f *FuncInstance) MethodLen() int {
	return len(f.methodMap)
}

// for sort

func (f *FuncInstance) Len() int {
	return len(f.keys)
}

func (f *FuncInstance) Swap(x int, y int) {
	old := f.keys[x]
	f.keys[x] = f.keys[y]
	f.keys[y] = old
}
func (f *FuncInstance) Less(x int, y int) bool {
	a := f.IndexOf(x)
	b := f.IndexOf(y)
	orderA := a.Order
	orderB := b.Order
	if orderA != "" && orderB == "" {
		return true
	}
	if orderA == "" && orderB != "" {
		return false
	}
	return orderA < orderB
}
