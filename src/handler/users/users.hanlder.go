package users

import (
	"evidence-hub-be/src/core/config"
	"evidence-hub-be/src/core/schema/common"
	"evidence-hub-be/src/core/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Tags Users
// @Summary Get All Users.
// @Description Retrieve all users entries.
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /v1/users [get]
// @securitydefinitions.apikey token
// @in header
// @name Authorization
// @Security token
func (h *Handler) GetAllUsers(ctx *gin.Context) error {

	var user common.Users
	var role common.Roles

	QroleName := strings.ToUpper(ctx.Query("role"))

	tokenString := ctx.GetHeader("Authorization")

	res, err := h.token.VerifyToken(tokenString)
	if err != nil {
		return utils.ResponseError("01", err.Error(), ctx)
	}

	if err := config.DB.
		Preload("Role").
		First(&user, "id = ?", res.UserID).Error; err != nil {

		return utils.ResponseError("99", "Unauthorized", ctx)
	}

	if user.Role.Name != "AUDIT" {
		return utils.ResponseError("99", "Forbidden", ctx)
	}
	var users []common.Users

	if QroleName != "" {

		if err := config.DB.First(&role, "name = ?", strings.ToUpper(QroleName)).Error; err != nil {
			return utils.ResponseError("99", "Role Not Found", ctx)
		}

		if err := config.DB.Preload("Role").Find(&users, "role_id = ?", role.ID).Error; err != nil {
			return utils.ResponseError("99", "User Not Found", ctx)
		}
	} else {
		if err := config.DB.Preload("Role").Find(&users).Error; err != nil {
			return utils.ResponseError("99", "User Not Found", ctx)
		}
	}

	return utils.ResponseSuccess("00", users, ctx)
}
