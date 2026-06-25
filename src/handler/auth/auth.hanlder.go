package auth

import (
	"evidence-hub-be/src/core/config"
	"evidence-hub-be/src/core/pkg/token"
	"evidence-hub-be/src/core/schema/common"
	"evidence-hub-be/src/core/utils"
	"evidence-hub-be/src/handler/auth/model"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// @Tags Authentication
// @Summary Register User.
// @Description Register User For PBB, Payload Role Must uuid of Ref Schema.
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /v1/register [post]
// @Param register body model.SPayloadRegister true "register"
func (h *Handler) SignUpHandler(ctx *gin.Context) error {
	var payload *model.SPayloadRegister
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		return utils.ResponseError("99", err.Error(), ctx)
	}

	roleSchema := new(common.Roles)
	if err := config.DB.First(roleSchema, "id = ?", payload.RoleID).Error; err != nil {
		return utils.ResponseError("99", "Invalid Role ID", ctx)
	}

	newUser := new(common.Users)
	payload.MapToSchema(newUser, &roleSchema.ID)

	if err := config.DB.Create(newUser).Error; err != nil {
		return utils.ResponseError("99", err.Error(), ctx)
	}

	return utils.ResponseSuccess("00", newUser, ctx)
}

// @Tags Authentication
// @Summary Login User.
// @Description Login User return token in body or if next js framework must withCredential false.
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /v1/login [post]
// @Param login body model.SPayloadLogin true "login"
// @securitydefinitions.apikey token
// @in header
// @name Authorization
// @Security token
func (h *Handler) SigninHandler(ctx *gin.Context) error {
	var payload *model.SPayloadLogin

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		return utils.ResponseError("99", err.Error(), ctx)
	}

	var user common.Users
	result := config.DB.Preload("Role").
		First(&user, "username = ?", strings.ToLower(payload.Username))
	if result.Error != nil {
		return utils.ResponseError("99", "Invalid Username or Password", ctx)
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return utils.ResponseError("99", "Invalid Username or Password", ctx)
	}

	tokenPayload := &token.Payload{
		UserID: user.ID.String(),
		Role: token.PayloadRole{
			ID:   user.RoleID.String(),
			Name: user.Role.Name,
		},
	}

	tokenString, err := h.token.GenerateToken(tokenPayload)
	if err != nil {
		return utils.ResponseError("99", fmt.Sprintf("generating JWT Token failed: %v", err), ctx)
	}

	jwt_maxage := h.cfg.AppConfig.JwtMaxAge
	maxage := jwt_maxage * 60
	ctx.SetCookie(
		"token",
		tokenString,
		int(maxage),
		"/",
		"",
		false,
		true,
	)
	resp := map[string]interface{}{
		"message": "Success Login Account",
		"data": map[string]interface{}{
			"token": tokenString,
			"user":  user,
		},
	}
	return utils.ResponseSuccess("00", resp, ctx)
}

// @Tags Authentication
// @Summary Token Check User.
// @Description Token Check User return token in body or if next js framework must withCredential false.
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /v1/token-check [get]
// @securitydefinitions.apikey token
// @in header
// @name Authorization
// @Security token
func (h *Handler) TokenCheck(ctx *gin.Context) error {
	var user common.Users
	tokenString := ctx.GetHeader("Authorization")
	res, err := h.token.VerifyToken(tokenString)
	if err != nil {
		return utils.ResponseError("01", err.Error(), ctx)
	}
	if err := config.DB.Preload("Role").First(&user, "id = ?", res.UserID).Error; err != nil {
		return utils.ResponseError("99", "Unauthorized", ctx)
	}
	userResponse := model.UserResponse{Users: user}
	return utils.ResponseSuccess("00", userResponse, ctx)
}

// @Tags Authentication
// @Summary Change Password.
// @Description Change Password.
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /v1/change-password [post]
// @Param change-password body model.SPayloadChangePassword true "Change Password"
// @securitydefinitions.apikey token
// @in header
// @name Authorization
// @Security token
func (h *Handler) ChangePassword(ctx *gin.Context) error {
	var payload *model.SPayloadChangePassword

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		return utils.ResponseError("99", err.Error(), ctx)
	}

	var user common.Users

	err := config.DB.
		Preload("Role").
		First(&user, "username = ?", strings.ToLower(payload.Username)).Error
	if err != nil {
		return utils.ResponseError("99", err.Error(), ctx)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.OldPassword))
	if err != nil {
		return utils.ResponseError("99", "Invalid Old Password", ctx)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return utils.ResponseError("99", err.Error(), ctx)
	}

	user.Password = string(hashedPassword)
	user.IsTempPassword = false
	user.TempPasswordExpiresAt = nil

	if err := config.DB.Save(&user).Error; err != nil {
		return utils.ResponseError("99", err.Error(), ctx)
	}

	userResponse := model.UserResponse{Users: user}
	return utils.ResponseSuccess("00", userResponse, ctx)
}
