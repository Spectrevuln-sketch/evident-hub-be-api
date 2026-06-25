package leaderboard

import (
	"evidence-hub-be/src/core/config"
	"evidence-hub-be/src/core/schema/common"
	"evidence-hub-be/src/core/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Tags Leaderboard
// @Summary Get All Leaderboard.
// @Description Retrieve all leaderboard entries.
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /v1/leaderboard [get]
// @securitydefinitions.apikey token
// @in header
// @name Authorization
// @Security token
func (h *Handler) GetAllLeaderBoard(ctx *gin.Context) error {

	var user common.Users

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

	type LeaderBoardResponse struct {
		Rank       int64     `json:"rank"`
		UserID     uuid.UUID `json:"user_id"`
		Fullname   string    `json:"fullname"`
		Username   string    `json:"username"`
		TotalPoint int64     `json:"total_point"`
		TotalFind  int64     `json:"total_find"`
	}

	var result []LeaderBoardResponse

	err = config.DB.Raw(`
		SELECT
			u.id AS user_id,
			u.fullname,
			u.username,
			COALESCE(SUM(lb.point),0) AS total_point,
			COUNT(lb.id) AS total_find
		FROM leader_boards lb
		JOIN users u
			ON u.id = lb.user_id
		GROUP BY
			u.id,
			u.fullname,
			u.username
		ORDER BY total_point DESC
	`).Scan(&result).Error

	if err != nil {
		return utils.ResponseError("99", err.Error(), ctx)
	}

	for i := range result {
		result[i].Rank = int64(i + 1)
	}

	return utils.ResponseSuccess("00", result, ctx)
}

// @Tags Leaderboard
// @Summary Get By Id Leaderboard.
// @Description Retrieve one leaderboard entries.
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /v1/leaderboard/{id} [get]
// @securitydefinitions.apikey token
// @in header
// @name Authorization
// @Security token
func (h *Handler) GetLeaderBoardById(ctx *gin.Context) error {

	id := ctx.Param("id")

	type LeaderBoardResponse struct {
		Rank       int64     `json:"rank"`
		UserID     uuid.UUID `json:"user_id"`
		Fullname   string    `json:"fullname"`
		Username   string    `json:"username"`
		TotalPoint int64     `json:"total_point"`
		TotalFind  int64     `json:"total_find"`
	}

	var result []LeaderBoardResponse

	err := config.DB.Raw(`
		SELECT
			u.id AS user_id,
			u.fullname,
			u.username,
			COALESCE(SUM(lb.point),0) AS total_point,
			COUNT(lb.id) AS total_find
		FROM leader_boards lb
		JOIN users u
			ON u.id = lb.user_id
		GROUP BY
			u.id,
			u.fullname,
			u.username
		ORDER BY total_point DESC
	`).Scan(&result).Error

	if err != nil {
		return utils.ResponseError("99", err.Error(), ctx)
	}

	for i := range result {

		result[i].Rank = int64(i + 1)

		if result[i].UserID.String() == id {
			return utils.ResponseSuccess("00", result[i], ctx)
		}
	}

	return utils.ResponseError("99", "Leaderboard not found", ctx)
}

// @Tags Leaderboard
// @Summary Get My Leaderboard.
// @Description Retrieve one leaderboard entries.
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /v1/leaderboard/me [get]
// @securitydefinitions.apikey token
// @in header
// @name Authorization
// @Security token
func (h *Handler) GetMyLeaderBoard(ctx *gin.Context) error {

	tokenString := ctx.GetHeader("Authorization")

	res, err := h.token.VerifyToken(tokenString)
	if err != nil {
		return utils.ResponseError("01", err.Error(), ctx)
	}

	type LeaderBoardResponse struct {
		Rank       int64     `json:"rank"`
		UserID     uuid.UUID `json:"user_id"`
		Fullname   string    `json:"fullname"`
		Username   string    `json:"username"`
		TotalPoint int64     `json:"total_point"`
		TotalFind  int64     `json:"total_find"`
	}

	var result []LeaderBoardResponse

	err = config.DB.Raw(`
		SELECT
			u.id AS user_id,
			u.fullname,
			u.username,
			COALESCE(SUM(lb.point),0) AS total_point,
			COUNT(lb.id) AS total_find
		FROM leader_boards lb
		JOIN users u
			ON u.id = lb.user_id
		GROUP BY
			u.id,
			u.fullname,
			u.username
		ORDER BY total_point DESC
	`).Scan(&result).Error

	if err != nil {
		return utils.ResponseError("99", err.Error(), ctx)
	}

	for i := range result {

		result[i].Rank = int64(i + 1)

		if result[i].UserID.String() == res.UserID {
			return utils.ResponseSuccess("00", result[i], ctx)
		}
	}

	return utils.ResponseError("99", "Leaderboard not found", ctx)
}
