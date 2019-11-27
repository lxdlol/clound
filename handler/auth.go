package handler

import (
	"github.com/gin-gonic/gin"
	"myyun/util"
	"net/http"
)

func HTTPInterceptor()gin.HandlerFunc {
	return func(c *gin.Context) {
			username := c.Request.FormValue("username")
			token :=c.Request.FormValue("token")

			//验证登录token是否有效
			if len(username) < 3 || !IsTokenValid(token) {
				c.Abort()
				resp:=util.NewRespMsg(
					-3,
					"token not use",
					nil,
					)
				c.JSON(http.StatusOK,resp.JSONBytes())
				// w.WriteHeader(http.StatusForbidden)
				// token校验失败则跳转到登录页面
				return
			}
			c.Next()
		}

}

