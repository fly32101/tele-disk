package middleware

import (
	"mime"
	"net/http"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Response 统一成功/失败返回结构。
// 与默认 MiddlewareHandlerResponse 类似，但输出固定字段 {code,message,data}，并避免覆盖流式响应。
func Response() ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		r.Middleware.Next()

		// 已有输出/流式响应则跳过包装。
		if r.Response.BufferLength() > 0 || r.Response.BytesWritten() > 0 {
			return
		}
		mediaType, _, _ := mime.ParseMediaType(r.Response.Header().Get("Content-Type"))
		for _, ct := range []string{"text/event-stream", "application/octet-stream", "multipart/x-mixed-replace"} {
			if mediaType == ct {
				return
			}
		}

		var (
			err   = r.GetError()
			data  = r.GetHandlerResponse()
			code  = gcode.CodeOK.Code()
			msg   = gcode.CodeOK.Message()
		)

		// 非 200 的状态码也视为错误
		if r.Response.Status > 0 && r.Response.Status != http.StatusOK && err == nil {
			switch r.Response.Status {
			case http.StatusNotFound:
				err = gerror.NewCode(gcode.CodeNotFound, "not found")
			case http.StatusForbidden:
				err = gerror.NewCode(gcode.CodeNotAuthorized, "forbidden")
			default:
				err = gerror.NewCode(gcode.CodeUnknown, http.StatusText(r.Response.Status))
			}
		}

		if err != nil {
			if c := gerror.Code(err); c != gcode.CodeNil {
				code = c.Code()
				msg = c.Message()
			} else {
				code = gcode.CodeInternalError.Code()
				msg = err.Error()
			}
		}

		r.Response.WriteJson(g.Map{
			"code":    code,
			"message": msg,
			"data":    data,
		})
	}
}
