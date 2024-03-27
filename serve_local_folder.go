package s2d3

import (
	"context"
	"net/http"

	"github.com/usalko/s2d3/services"
)

type ServeLocalFolder struct {
	RootFolder string `default:"."`
	UrlContext string `default:""`
	ServerAddr string `default:"localhost:8081"`
}

func (serveLocalFolder *ServeLocalFolder) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	ctx = context.WithValue(ctx, services.KeyServerAddr, serveLocalFolder.ServerAddr)
	ctx = context.WithValue(ctx, services.KeyDataFolder, serveLocalFolder.RootFolder)
	ctx = context.WithValue(ctx, services.KeyUrlContext, serveLocalFolder.UrlContext)
	services.ApiRouter(writer, request.WithContext(ctx))
}
