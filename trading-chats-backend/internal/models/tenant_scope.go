package models

func IsAdmin(ctx *AuthContext) bool {
	return ctx != nil && ctx.Role == RoleAdmin
}

func ResolveTenantID(ctx *AuthContext, explicitTenantID string) string {
	if explicitTenantID != "" {
		return explicitTenantID
	}
	if ctx == nil {
		return ""
	}
	return ctx.TenantID
}
