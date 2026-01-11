package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Stat(c *gin.Context) {
	url := c.Request.URL.Path
	fmt.Printf("request url is: ", url)
}
