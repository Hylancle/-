package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
	"xiaoweishu/webook/internal/web"
)

// Gin的middleware可以做到非常多时期，包括熔断限流降级、日志 metrics。身份认证和健全
func main() {
	hdl := web.NewUserHandler()
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		// 允许自带信息，cookie
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"X-Jwt-Token", "X-Refresh-Token"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return false
		},
		MaxAge: 12 * time.Hour,
	}))
	hdl.RegisterRoutes(server)
	server.Run()

}
