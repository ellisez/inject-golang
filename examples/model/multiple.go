package model

type Event struct {
	Config *Config
	*Database
	EventName string
	Data      map[string]any
}

type Listener struct {
	EventName string
	Func      func(map[string]any)
}
