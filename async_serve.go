/**
 * Author: Vanya Usalko <ivict@rambler.ru>
 * File: s2d3.go
 */

package s2d3

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/usalko/s2d3/services"
)

func InitStorage(localFolder string) {
	storage := services.Storage{
		RootFolder: localFolder,
	}
	storage.Init()
}

func AsyncServe(localFolder string, addr string, port int) (context.Context, context.CancelFunc) {
	fmt.Printf("Serve local folder '%s' \n", localFolder)
	fmt.Printf("Host: %s Port: %d \n", addr, port)

	multiplexer := http.NewServeMux()
	multiplexer.HandleFunc(".*", services.ApiRouter)
	// multiplexer.HandleFunc("/hello", services.GetHello)

	ctx, cancelFunc := context.WithCancel(context.Background())
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", addr, port),
		Handler: multiplexer,
		BaseContext: func(listener net.Listener) context.Context {
			ctx = context.WithValue(ctx, services.KeyServerAddr, listener.Addr().String())
			ctx = context.WithValue(ctx, services.KeyDataFolder, localFolder)
			ctx = context.WithValue(ctx, services.KeyStatisticsApplicationFolder, os.Getenv("STATISTICS_APPLICATION_FOLDER"))
			return ctx
		},
	}
	go func() {
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server closed\n")
		} else if err != nil {
			fmt.Printf("error listening for server: %s\n", err)
		}
		cancelFunc()
	}()

	return ctx, cancelFunc
}
