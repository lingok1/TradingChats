package models

import "context"

type contextKey string

const authContextStorageKey contextKey = "auth_context"

func WithAuthContext(ctx context.Context, auth *AuthContext) context.Context {
	return context.WithValue(ctx, authContextStorageKey, auth)
}

func GetAuthContext(ctx context.Context) *AuthContext {
	if ctx == nil {
		return &AuthContext{}
	}
	value := ctx.Value(authContextStorageKey)
	authCtx, _ := value.(*AuthContext)
	if authCtx == nil {
		return &AuthContext{}
	}
	return authCtx
}
