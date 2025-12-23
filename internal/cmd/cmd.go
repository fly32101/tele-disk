package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"tele-disk/internal/controller/auth"
	"tele-disk/internal/controller/files"
	"tele-disk/internal/controller/hello"
	"tele-disk/internal/controller/telegram"
	"tele-disk/internal/logic/middleware"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			// Serve static front-end from resource/public
			s.SetServerRoot("resource/public")
			s.SetIndexFiles([]string{"html/index.html", "index.html"})
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.Response())
				group.Bind(
					auth.NewV1(),
					hello.NewV1(),
				)
				group.Group("/", func(group *ghttp.RouterGroup) {
					group.Middleware(middleware.Auth())
					group.Bind(
						files.NewV1(),
						telegram.NewV1(),
					)
				})
			})
			s.Run()
			return nil
		},
	}
)
