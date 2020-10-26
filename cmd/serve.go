package cmd

import (
	"fmt"
	kittyhttp "github.com/Reasno/kitty/pkg/http"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	"github.com/oklog/run"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Long:  `Start the gRPC server and HTTP server`,
	PreRunE: initServiceContainer,
	Run: func(cmd *cobra.Command, args []string) {
		var g run.Group

		// Run all exit logic
		defer func() {
			for _, f := range serviceContainer.CloserProviders {
				f()
			}
		}()

		// Start HTTP Server
		{
			httpAddr := viper.GetString("global.http.addr")
			ln, err := net.Listen("tcp", httpAddr)
			if err != nil {
				_ = logger.Log("err", err)
				os.Exit(1)
			}
			g.Add(func() error {
				h := getHttpHandler(ln, serviceContainer.HttpProviders...)
				return http.Serve(ln, h)
			}, func(err error) {
				_ = ln.Close()
			})
		}

		// Start gRPC server
		{
			grpcAddr := viper.GetString("global.grpc.addr")
			ln, err := net.Listen("tcp", grpcAddr)
			if err != nil {
				_ = logger.Log("err", err)
				os.Exit(1)
			}
			g.Add(func() error {
				s := getGRPCServer(ln, serviceContainer.GrpcProviders...)
				return s.Serve(ln)
			}, func(err error) {
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
		for _, s := range serviceContainer.RunProviders {
			g.Add(s.Loop, s.Exit)
		}

		// Add Cronjob etc. here

		if err := g.Run(); err != nil {
			level.Warn(logger).Log("err", err)
		}

		level.Info(logger).Log("msg", "graceful shutdown complete; see you next time")
	},
}

func getHttpHandler(ln net.Listener, providers ...func(*mux.Router)) http.Handler {
	_ = logger.Log("transport", "HTTP", "addr", ln.Addr())

	var handler http.Handler
	var router = mux.NewRouter()
	for _, p := range providers {
		p(router)
	}
	handler = kittyhttp.AddCorsMiddleware()(router)
	return handler
}

func getGRPCServer(ln net.Listener, providers ...func(s *grpc.Server)) *grpc.Server {
	_ = logger.Log("transport", "gRPC", "addr", ln.Addr())

	s := grpc.NewServer()
	for _, p := range providers {
		p(s)
	}
	return s
}
