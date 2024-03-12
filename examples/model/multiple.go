package model

import "fmt"

type MiddleWare struct {
	config *Config
	*Database
	Path   string
	Handle func() error
}

func (m *MiddleWare) Config() *Config {
	fmt.Println("call MiddleWare.config getter")
	return m.config
}

func (m *MiddleWare) SetConfig(config *Config) {
	fmt.Println("call MiddleWare.config setter")
	m.config = config
}

type Router struct {
	MiddleWare *MiddleWare
	Path       string
	Handle     func() error
}
