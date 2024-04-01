package main

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"io"
	"log"
	"net"
	"net/http"
)

type EchoHandler struct {
}

func (e *EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := io.Copy(w, r.Body); err != nil {
		log.Println("error occurred")
	}

}

func NewEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func NewServeMux(echoHandler *EchoHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/echo", echoHandler)
	return mux
}

func NewHTTPServer(lc fx.Lifecycle, mux *http.ServeMux) *http.Server {
	srv := &http.Server{
		Addr:    ":8000",
		Handler: mux,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			fmt.Println("Starting HTTP Server at", srv.Addr)
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}

func main() {
	fx.New(
		fx.Provide(NewHTTPServer),
		fx.Provide(NewServeMux),
		fx.Provide(NewEchoHandler),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}
