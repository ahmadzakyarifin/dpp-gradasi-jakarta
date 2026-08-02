package middleware

import "github.com/gin-gonic/gin"

func GetUserID(c *gin.Context) (uint, bool) {
	v, ok := c.Get(ContextUserID)
	if !ok {
		return 0, false
	}

	id, ok := v.(uint)
	return id, ok
}

func GetRoleID(c *gin.Context) (uint, bool) {
	v, ok := c.Get(ContextRoleID)
	if !ok {
		return 0, false
	}

	id, ok := v.(uint)
	return id, ok
}

func GetEmail(c *gin.Context) (string, bool) {
	v, ok := c.Get(ContextEmail)
	if !ok {
		return "", false
	}

	email, ok := v.(string)
	return email, ok
}

func GetIPAddress(c *gin.Context) string {
	return c.ClientIP()
}

func GetUserAgent(c *gin.Context) string {
	return c.Request.UserAgent()
}
