package ginIOC_param

import (
	"github.com/gin-gonic/gin"
)

// 读取form值
func FormValue(c *gin.Context, key string) string {
	req := c.Request
	if req.Method == "POST" {
		req.ParseForm()
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req.FormValue(key)
}

func VarsValue(c *gin.Context, key string) string {
	return c.Param(key)
}