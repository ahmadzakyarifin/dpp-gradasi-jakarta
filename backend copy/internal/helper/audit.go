package helper

import (
	"context"
)

func GetAuditMeta(ctx context.Context) (uint, string, string, string, string) {
	if ctx == nil {
		return 0, "", "", "", ""
	}

	userID, _ := ctx.Value("user_id").(uint)
	userName, _ := ctx.Value("user_name").(string)
	role, _ := ctx.Value("role_name").(string)
	ipAddress, _ := ctx.Value("ip_address").(string)
	userAgent, _ := ctx.Value("user_agent").(string)

	return userID, userName, role, ipAddress, userAgent
}
