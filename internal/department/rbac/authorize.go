package rbac

import (
	"app/internal/http/middleware"
	"fmt"
)

const adminRole = "admin"

func Authorize(ctx *middleware.AuthContext, scope string, errorMessage string) error {
	if isAdmin(ctx) {
		return nil
	}

	if hasScope(ctx, scope) {
		return nil
	}

	return fmt.Errorf(errorMessage)
}

func isAdmin(ctx *middleware.AuthContext) bool {
	for _, role := range ctx.Roles {
		if role == adminRole {
			return true
		}
	}

	return false
}

func hasScope(ctx *middleware.AuthContext, scope string) bool {
	for _, s := range ctx.Scopes {
		if s == scope {
			return true
		}
	}

	return false
}
