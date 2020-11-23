package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	"github.com/oklog/run"
	"github.com/spf13/cobra"
	kittyhttp "glab.tagtic.cn/ad_gains/kitty/pkg/khttp"
	"google.golang.org/grpc"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Long:  `Start the gRPC server and HTTP server`,
	Run: func(cmd *cobra.Command, args []string) {
		initModules()
		defer shutdownModules()

		var g run.Group

		// Start HTTP Server
		{
			httpAddr := conf().String("global.http.addr")
			ln, err := net.Listen("tcp", httpAddr)
			if err != nil {
				er(err)
				os.Exit(1)
			}
			h := getHttpHandler(ln, moduleContainer.HttpProviders...)
			srv := &http.Server{
				Handler:      h,
				IdleTimeout:  2 * time.Second,
				ReadTimeout:  2 * time.Second,
				WriteTimeout: 2 * time.Second,
			}
			g.Add(func() error {
				return srv.Serve(ln)
			}, func(err error) {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := srv.Shutdown(ctx); err != nil {
					_ = level.Warn(coreModule.Logger).Log("err", err)
					os.Exit(1)
				}
				_ = ln.Close()

			})
		}

		// Start gRPC server
		{
			grpcAddr := conf().String("global.grpc.addr")
			ln, err := net.Listen("tcp", grpcAddr)
			if err != nil {
				_ = level.Error(coreModule.Logger).Log("err", err)
				os.Exit(1)
			}
			s := getGRPCServer(ln, moduleContainer.GrpcProviders...)
			g.Add(func() error {
				return s.Serve(ln)
			}, func(err error) {
				s.GracefulStop()
				_ = ln.Close()
			})
		}

		// Graceful shutdown
		{
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			g.Add(func() error {
				terminateError := fmt.Errorf("%s", <-c)
				return terminateError
			}, func(err error) {
				close(c)
			})
		}

		// Additional run groups
		for _, s := range moduleContainer.RunProviders {
			s(&g)
		}

		// Add Cronjob etc. here

		if err := g.Run(); err != nil {
			er(err)
		}

		info("graceful shutdown complete; see you next time")
	},
}

func getHttpHandler(ln net.Listener, providers ...func(*mux.Router)) http.Handler {
	_ = level.Info(coreModule.Logger).Log("transport", "HTTP", "addr", ln.Addr())

	var handler http.Handler
	var router = mux.NewRouter()
	for _, p := range providers {
		p(router)
	}
	handler = kittyhttp.AddCorsMiddleware()(router)
	return handler
}

func getGRPCServer(ln net.Listener, providers ...func(s *grpc.Server)) *grpc.Server {
	_ = level.Info(coreModule.Logger).Log("transport", "gRPC", "addr", ln.Addr())

	s := grpc.NewServer(grpc.ConnectionTimeout(time.Second))
	for _, p := range providers {
		p(s)
	}
	return s
}
