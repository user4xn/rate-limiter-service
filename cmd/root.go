package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"os"
	"rate-limiter/pkg/factory"
	"rate-limiter/pkg/route"
	"rate-limiter/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "start rate limiter service",
	Run: func(cmd *cobra.Command, args []string) {
		f := factory.NewFactory()
		g := gin.New()

		route.NewAPIHttp(g, f)

		svc := http.Server{
			Addr:    ":" + util.GetEnv("APP_PORT", "8080"),
			Handler: g,
		}

		log.Printf("start listening on %s", util.GetEnv("APP_PORT", "8080"))

		go func() {
			if err := svc.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("failed to start http server: %s", err.Error())
			}
		}()

		wait := util.GracefulShutdown(context.Background(), time.Second*10, map[string]func(ctx context.Context) error{
			"http-server": func(ctx context.Context) error {
				return svc.Shutdown(ctx)
			},
		})

		<-wait
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
