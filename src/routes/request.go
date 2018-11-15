package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http/httputil"
	"strconv"
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
	return c.Query(key)
}

func VarsValueInt(c *gin.Context, key string) (int, error) {
	value := c.Query(key)
	num, err := strconv.ParseInt(value, 10, 32)
	return int(num), err
}

func VarsValueFloat(c *gin.Context, key string, bitsize int) (float64, error) {
	value := c.Query(key)
	return strconv.ParseFloat(value, bitsize)
}

func DumpRequest(c *gin.Context) {
	data, err := httputil.DumpRequest(c.Request, true)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(data))
	}
}
