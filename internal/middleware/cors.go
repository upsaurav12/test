package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS returns a CORS middleware configured with the supplied allowed origins.
// In production pass an explicit allow-list; avoid wildcard "*" for credentialed
// requests.
func CORS(allowedOrigins []string) gin.HandlerFunc {
	cfg := cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", RequestIDHeader},
		ExposeHeaders:    []string{RequestIDHeader},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
	return cors.New(cfg)
}
