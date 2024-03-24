/**
 * Author: Vanya Usalko <ivict@rambler.ru>
 * File: s2d3.go
 */

package s2d3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
)

type KeyServerAddr string

const keyServerAddr KeyServerAddr = "serverAddr"

func getRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("%s: got / request\n", ctx.Value(keyServerAddr))
	io.WriteString(w, "This is my website!\n")
}

func getHello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("%s: got /hello request\n", ctx.Value(keyServerAddr))
	io.WriteString(w, "Hello, HTTP!\n")
}

func Serve(localFolder string) (context.Context, context.CancelFunc) {
	println("Go on")

	multiplexer := http.NewServeMux()
	multiplexer.HandleFunc("/", getRoot)
	multiplexer.HandleFunc("/hello", getHello)

	ctx, cancelFunc := context.WithCancel(context.Background())
	server := &http.Server{
		Addr:    ":3333",
		Handler: multiplexer,
		BaseContext: func(listener net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, listener.Addr().String())
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
