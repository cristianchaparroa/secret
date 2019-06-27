package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func removeParams(c *gin.Context) string {
	url := c.Request.URL.String()

	fmt.Println(c.Params)
	for _, p := range c.Params {
		url = strings.Replace(url, p.Value, "", 1)
	}
	return url
}
