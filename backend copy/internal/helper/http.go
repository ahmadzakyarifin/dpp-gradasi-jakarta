package helper

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
)

// RestoreBody mengembalikan request body yang sebelumnya
// telah disimpan Gin melalui ShouldBindBodyWith.
func RestoreBody(c *gin.Context) {
	value, ok := c.Get(gin.BodyBytesKey)
	if !ok {
		return
	}

	body, ok := value.([]byte)
	if !ok {
		return
	}

	c.Request.Body = io.NopCloser(
		bytes.NewReader(body),
	)
}
