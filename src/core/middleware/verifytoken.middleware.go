package middleware

import (
	pkgctx "evidence-hub-be/src/core/pkg/context"
	"evidence-hub-be/src/core/utils"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) VerifyTokenMiddleware(ctx *gin.Context) error {
	token := ctx.GetHeader("Authorization")
	res, err := m.token.VerifyToken(token)
	if err != nil {
		return utils.ResponseError("01", err.Error(), ctx)
	}

	ctx.Set(pkgctx.UserIDKey, res.UserID)
	ctx.Set(pkgctx.RoleIDKey, res.RoleID)
	ctx.Set(pkgctx.RoleKey, res.Role)

	ctx.Next()
	return nil
}
