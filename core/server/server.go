package server

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tylerb/graceful"
)

// Server server
type Server struct {
	e    *echo.Echo
	port string
}

// New server
func New(router *echo.Echo, port string) *Server {
	return &Server{
		e:    router,
		port: port,
	}
}

// Start start server
func (s *Server) Start() {
	s.e.Server.Addr = s.port
	s.e.Logger.Infof("http server started on %s", s.e.Server.Addr)
	_ = graceful.ListenAndServe(s.e.Server, 5*time.Second)
}
