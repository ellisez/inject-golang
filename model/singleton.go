package model

type SingletonInstance struct {
	keys []*Key // singleton | argument | web

	provideMap map[string]*Provide

	argumentMap map[string]*Provide

	webMap map[string]*Provide

	webApplicationMap map[string]*WebApplication
}

func newCtxSingletonInstance() *SingletonInstance {
	return &SingletonInstance{
		provideMap:        map[string]*Provide{},
		argumentMap:       map[string]*Provide{},
		webMap:            map[string]*Provide{},
		webApplicationMap: map[string]*WebApplication{},
	}
}

func (s *SingletonInstance) Add(provide *Provide) {
	switch provide.Mode {
	case "singleton":
		s.keys = append(s.keys, &Key{
			Order:    provide.Order,
			Instance: provide.Instance,
			Type:     provide.Mode,
		})
		s.provideMap[provide.Instance] = provide
	case "argument":
		s.keys = append(s.keys, &Key{
			Order:    provide.Order,
			Instance: provide.Instance,
			Type:     provide.Mode,
		})
		s.argumentMap[provide.Instance] = provide
	}

}
func (s *SingletonInstance) AddWeb(provide *Provide, webApplication *WebApplication) {
	s.keys = append(s.keys, &Key{
		Order:    provide.Order,
		Instance: provide.Instance,
		Type:     "web",
	})
	s.webMap[provide.Instance] = provide

	s.webApplicationMap[provide.Instance] = webApplication

}
func (s *SingletonInstance) Delete(key string) {
	var hitKey *Key
	for i, k := range s.keys {
		if k.Instance == key {
			s.keys = append(s.keys[:i], s.keys[i+1:]...)
			hitKey = k
			break
		}
	}
	if hitKey != nil {
		switch hitKey.Type {
		case "singleton":
			delete(s.provideMap, key)
		case "argument":
			delete(s.argumentMap, key)
		case "web":
			delete(s.webMap, key)
			delete(s.webApplicationMap, key)
		}
	}
}

func (s *SingletonInstance) Contains(key string) bool {
	return s.KeyOf(key) != nil
}

func (s *SingletonInstance) GetProvide(key string) *Provide {
	return s.provideMap[key]
}

func (s *SingletonInstance) GetArgument(key string) *Provide {
	return s.argumentMap[key]
}

func (s *SingletonInstance) GetWebApplication(key string) (*Provide, *WebApplication) {
	return s.webMap[key], s.webApplicationMap[key]
}

func (s *SingletonInstance) Get(key string) (*Provide, *WebApplication) {
	k := s.KeyOf(key)
	if k == nil {
		return nil, nil
	}
	switch k.Type {
	case "singleton":
		return s.provideMap[key], nil
	case "argument":
		return s.argumentMap[key], nil
	case "web":
		return s.webMap[key], s.webApplicationMap[key]
	}
	return nil, nil
}

func (s *SingletonInstance) Replace(provide *Provide) bool {
	key := provide.Instance
	k := s.KeyOf(key)
	if k == nil {
		return false
	}
	if k.Type != provide.Mode {
		switch k.Type {
		case "singleton":
			delete(s.provideMap, key)
		case "argument":
			delete(s.argumentMap, key)
		case "web":
			delete(s.webMap, key)
			delete(s.webApplicationMap, key)
		}
		k.Type = provide.Mode
	}
	k.Order = provide.Order
	switch provide.Mode {
	case "singleton":
		s.provideMap[key] = provide
	case "argument":
		s.argumentMap[key] = provide
	default:
		return false
	}
	return true
}

func (s *SingletonInstance) ReplaceWeb(provide *Provide, webApplication *WebApplication) bool {
	key := provide.Instance
	k := s.KeyOf(key)
	if k == nil {
		return false
	}
	if k.Type != "web" {
		switch k.Type {
		case "singleton":
			delete(s.provideMap, key)
		case "argument":
			delete(s.argumentMap, key)
		case "web":
			delete(s.webMap, key)
			delete(s.webApplicationMap, key)
		}
		k.Type = "web"
	}

	k.Order = provide.Order
	s.webMap[provide.Instance] = provide
	s.webApplicationMap[provide.Instance] = webApplication
	return true
}

func (s *SingletonInstance) IndexOf(index int) (*Provide, *WebApplication) {
	key := s.keys[index]
	return s.Get(key.Instance)
}

func (s *SingletonInstance) KeyOf(key string) *Key {
	for _, k := range s.keys {
		if k.Instance == key {
			return k
		}
	}
	return nil
}
func (s *SingletonInstance) ProvideLen() int {
	return len(s.provideMap)
}
func (s *SingletonInstance) ArgumentLen() int {
	return len(s.argumentMap)
}
func (s *SingletonInstance) WebLen() int {
	return len(s.webMap)
}

// for sort

func (s *SingletonInstance) Len() int {
	return len(s.keys)
}

func (s *SingletonInstance) Swap(x int, y int) {
	old := s.keys[x]
	s.keys[x] = s.keys[y]
	s.keys[y] = old
}
func (s *SingletonInstance) Less(x int, y int) bool {
	a := s.keys[x]
	b := s.keys[y]
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
