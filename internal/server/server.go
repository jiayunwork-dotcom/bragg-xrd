package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Handler http.Handler
	Addr    string
}

func New(addr string) *Server {
	if addr == "" {
		addr = ":8080"
	}
	return &Server{Handler: Routes(), Addr: addr}
}

func (s *Server) ListenAndServe() error {
	httpServer := &http.Server{Addr: s.Addr, Handler: s.Handler,
		ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second,
		WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second}
	return httpServer.ListenAndServe()
}

func (s *Server) RunWithGracefulShutdown() error {
	httpServer := &http.Server{Addr: s.Addr, Handler: s.Handler,
		ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second,
		WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second}
	errCh := make(chan error, 1)
	go func() { errCh <- httpServer.ListenAndServe() }()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <-errCh:
		return err
	case <-sig:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return httpServer.Shutdown(ctx)
	}
}

func IsServerClosed(err error) bool {
	return errors.Is(err, http.ErrServerClosed)
}
