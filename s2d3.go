/**
 * Author: Vanya Usalko <ivict@rambler.ru>
 * File: s2d3.go
 */

package s2d3

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"

	"github.com/usalko/s2d3/services"
)

func InitStorage(localFolder string) {
	os.Mkdir(localFolder, fs.ModeAppend)
}

func Serve(localFolder string) (context.Context, context.CancelFunc) {
	println("Go on")

	multiplexer := http.NewServeMux()
	multiplexer.HandleFunc(".*", services.GetRoot)
	// multiplexer.HandleFunc("/hello", services.GetHello)

	ctx, cancelFunc := context.WithCancel(context.Background())
	server := &http.Server{
		Addr:    ":3333",
		Handler: multiplexer,
		BaseContext: func(listener net.Listener) context.Context {
			ctx = context.WithValue(ctx, services.KeyServerAddr, listener.Addr().String())
			ctx = context.WithValue(ctx, services.KeyDataFolder, localFolder)
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
