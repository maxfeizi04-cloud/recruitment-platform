package middleware

import "github.com/gin-gonic/gin"

func CharsetUTF8() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.Next()
	}
}
