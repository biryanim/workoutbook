package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

func XSSMiddleware() gin.HandlerFunc {
	p := bluemonday.UGCPolicy()

	return func(c *gin.Context) {
		for _, values := range c.Request.URL.Query() {
			for i, value := range values {
				values[i] = p.Sanitize(value)
			}
		}

		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			if err := c.Request.ParseForm(); err == nil {
				for _, values := range c.Request.PostForm {
					for i, value := range values {
						values[i] = p.Sanitize(value)
					}
				}
			}
		}

		c.Next()
	}
}
