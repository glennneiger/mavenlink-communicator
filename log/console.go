package log

import (
	"github.com/micro/go-log"
	"github.com/micro/go-micro/server"
	"golang.org/x/net/context"
	"time"
)

// ConsoleLogWrapper implements the server.HandlerWrapper for
// logging information to the console
func ConsoleLogWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		log.Logf("[%v] server request: %s\n", time.Now(), req.Method())
		return fn(ctx, req, rsp)
	}
}
