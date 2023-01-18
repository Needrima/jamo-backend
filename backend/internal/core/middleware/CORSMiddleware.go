package middleware

import (
	"github.com/gin-contrib/cors"
	"time"
)

var CORSMiddleware = cors.New(cors.Config{
	AllowOrigins:     []string{"*"},
	AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
	AllowHeaders:     []string{"Origin", "Content-Type"},
	ExposeHeaders:    []string{"Content-Length", "Content-Type"},
	AllowCredentials: true,
	AllowOriginFunc: func(origin string) bool {
		return true
	},
	MaxAge: 12 * time.Hour,
})
