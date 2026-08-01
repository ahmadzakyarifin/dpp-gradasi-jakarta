package helper

import (
	"context"
)

func GetAuditMeta(ctx context.Context) (uint, string, string, string, string) {
	if ctx == nil {
		return 0, "", "", "", ""
	}

	userID, _ := ctx.Value(ContextUserID).(uint)
	userName, _ := ctx.Value(ContextUserName).(string)
	role, _ := ctx.Value(ContextRoleName).(string)
	ipAddress, _ := ctx.Value(ContextIPAddress).(string)
	userAgent, _ := ctx.Value(ContextUserAgent).(string)

	return userID, userName, role, ipAddress, userAgent
}
