package cmd

import (
	"github.com/Reasno/kitty/app/handlers"
	"github.com/Reasno/kitty/app/svc"
	"github.com/Reasno/kitty/internal"
	pb "github.com/Reasno/kitty/proto"
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
		var g run.Group

		// Start HTTP Server
		{
			httpAddr := viper.GetString("http.addr")
			ln, err := net.Listen("tcp", httpAddr)
			if err != nil {
				_ = logger.Log("err", err)
				os.Exit(1)
			}
			g.Add(func() error {
				endpoints := internal.InitializeEndpoints()
				_ = logger.Log("transport", "HTTP", "addr", httpAddr)

				h := svc.MakeHTTPHandler(endpoints)

				router := mux.NewRouter()
				//router.Handle("/", http.StripPrefix("/", h))
				router.PathPrefix("/doc/").Handler(getOpenAPIHandler())
				router.PathPrefix("/").Handler(h)
				return http.Serve(ln,  router)
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
				endpoints := internal.InitializeEndpoints()
				_ = logger.Log("transport", "gRPC", "addr", grpcAddr)

				g := svc.MakeGRPCServer(endpoints)
				s := grpc.NewServer()
				pb.RegisterAppServer(s, g)
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

// getOpenAPIHandler serves an OpenAPI UI.
func getOpenAPIHandler() http.Handler {
	return http.StripPrefix("/doc",http.FileServer(http.Dir("./doc")))
}
