package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.New()
	r.GET("/berita/:slug", func(c *gin.Context) {})
	r.GET("/berita/admin", func(c *gin.Context) {})
}
