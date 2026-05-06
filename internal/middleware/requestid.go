package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// RequestIDKey is the gin context key for the request ID value.
	RequestIDKey = "request_id"
	// RequestIDHeader is the HTTP header name used to propagate request IDs.
	RequestIDHeader = "X-Request-ID"
)

// RequestID injects a unique request ID into each request's context and echoes
// it back in the response header. If the incoming request already carries an
// X-Request-ID header that value is reused.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(RequestIDHeader)
		if id == "" {
			id = uuid.New().String()
		}
		c.Set(RequestIDKey, id)
		c.Header(RequestIDHeader, id)
		c.Next()
	}
}
