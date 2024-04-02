package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/The-Flash/branch-lock/branch_map"
	"go.uber.org/fx"
)

func NewHTTPServer(lc fx.Lifecycle) *http.Server {
	srv := &http.Server{
		Addr: ":8000",
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
		fx.Provide(
			fx.Annotate(
				branch_map.NewBranchMap,
				fx.ResultTags(`name:"branchmap"`),
			),
		),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}
