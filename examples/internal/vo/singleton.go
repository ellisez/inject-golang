package vo

import (
	"fmt"
	"github.com/ellisez/inject-golang/examples/internal/entity"
)

type Config struct {
	Host string
	Port uint
}

type Database struct {
	Host     string
	Port     uint
	Schema   string
	UserName string
	Password string
}
type ServerInterface interface {
	Startup()
	Shutdown()
	IsRunning() bool
	TriggerEvent(eventName string, data map[string]any)
	AddListener(eventName string, handler func(map[string]any))

	EmptyListener()
}
type Server struct {
	Config    *Config
	Database  *Database
	Events    []*Event
	Listeners []*Listener
	idle      chan bool
}

func (s *Server) Startup() {
	if s.IsRunning() {
		return
	}
	s.idle = make(chan bool)
	go func() {
		for {
			for _, event := range s.Events {
				for _, handler := range s.Listeners {
					if event.EventName == handler.EventName {
						handler.Func(event.Data)
					}
				}
			}
			if !<-s.idle {
				break
			}
		}
		close(s.idle)
		s.idle = nil
	}()
}

func (s *Server) Shutdown() {
	if s.IsRunning() {
		s.idle <- false
	}
}

func (s *Server) IsRunning() bool {
	return s.idle != nil
}

func (s *Server) TriggerEvent(eventName string, data map[string]any) {
	e := &Event{
		EventName: eventName,
		Data:      data,
	}
	s.Events = append(s.Events, e)
	s.idle <- true
}

func (s *Server) AddListener(eventName string, handler func(map[string]any)) {
	s.Listeners = append(s.Listeners, &Listener{
		EventName: eventName,
		Func:      handler,
	})
}

func (s *Server) EmptyListener() {
	s.Listeners = nil
}

// TestServer example for inject method with uninjected recv
// @proxy
// @import github.com/ellisez/inject-golang/examples/internal/entity en
// @injectParam database Database
// @injectFunc FindAccountByName
// !@injectCtx appCtx
// @injectCall [ellisAccount] GetEllisAccount
func (s *Server) TestServer( /*appCtx ctx.Ctx, */ FindAccountByName func(string) *entity.Account, ellisAccount *entity.Account, database *Database) {
	fmt.Printf("call TestServer: %v, %v\n", s, database)
	s.AddListener("login", func(data map[string]any) {
		username := data["username"].(string)
		password := data["password"]
		account := FindAccountByName(username)
		if account == nil {
			fmt.Printf(`account "%s" is not found`, username)
			return
		}
		if account.Password != password {
			fmt.Printf(`account "%s" password is incorrect`, username)
			return
		}
		fmt.Printf(`account "%s" login succeeded`, username)
	})
	s.TriggerEvent("login", map[string]any{
		"username": "ellis",
		"password": "123456",
	})
	s.EmptyListener()
}
