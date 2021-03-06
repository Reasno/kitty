package cmd

import (
	"github.com/Reasno/kitty/app/handlers"
	"github.com/Reasno/kitty/app/register"
	kitty_http_middleware "github.com/Reasno/kitty/pkg/middleware/http"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	"github.com/oklog/run"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Long:  `Start the gRPC server and HTTP server`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			g         run.Group
			httpProviders []func() http.Handler
			grpcProviders []func(server2 *grpc.Server)
		)

		// Register generated services
		{
			register.RegisterApp(&httpProviders, &grpcProviders)
		}

		// Start HTTP Server
		{
			httpAddr := viper.GetString("http.addr")
			ln, err := net.Listen("tcp", httpAddr)
			if err != nil {
				_ = logger.Log("err", err)
				os.Exit(1)
			}
			g.Add(func() error {
				h := getHttpHandler(ln, httpProviders...)
				return http.Serve(ln, h)
			}, func(err error) {
				_ = ln.Close()
			})
		}

		// Start gRPC server
		{
			grpcAddr := viper.GetString("grpc.addr")
			ln, err := net.Listen("tcp", grpcAddr)
			if err != nil {
				_ = logger.Log("err", err)
				os.Exit(1)
			}
			g.Add(func() error {
				s := getGRPCServer(ln, grpcProviders...)
				return s.Serve(ln)
			}, func(err error) {
				_ = ln.Close()
			})
		}

		// Graceful shutdown
		{
			var ch chan error
			g.Add(func() error {
				ch = make(chan error)
				go handlers.InterruptHandler(ch)
				err := <-ch
				return err
			}, func(err error) {
				close(ch)
			})
		}

		// Add Cronjob etc. here

		if err := g.Run(); err != nil {
			level.Warn(logger).Log("err", err)
		}

		level.Info(logger).Log("msg", "graceful shutdown complete; see you next time")
	},
}

func getHttpHandler(ln net.Listener, providers... func() http.Handler) http.Handler {
	_ = logger.Log("transport", "HTTP", "addr", ln.Addr())

	var r http.Handler
	for _, p := range providers {
		rr := p().(*mux.Router)
		rr.Handle("/", r)
		r = rr
	}
	r = kitty_http_middleware.AddMetricMiddleware()(r)
	r = kitty_http_middleware.AddDocMiddleware()(r)
	r = kitty_http_middleware.AddCorsMiddleware()(r)
	r = kitty_http_middleware.AddHealthCheck()(r)
	return r
}

func getGRPCServer(ln net.Listener, providers... func(s *grpc.Server)) *grpc.Server {
	_ = logger.Log("transport", "gRPC", "addr", ln.Addr())

	s := grpc.NewServer()
	for _, p := range providers {
		p(s)
	}
	return s
}

