package evident

import (
	"evidence-hub-be/src/core/config"
	"evidence-hub-be/src/core/schema/common"
	"evidence-hub-be/src/core/schema/evident"
	"evidence-hub-be/src/core/schema/evidentPhoto"
	"evidence-hub-be/src/core/schema/leaderboard"
	"evidence-hub-be/src/core/utils"
	"evidence-hub-be/src/handler/evident/model"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Tags Evident
// @Summary Create New Evident
// @Description Create a new evident entry
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Title"
// @Param description formData string true "Description"
// @Param upload_temuan formData file false "Foto Evident"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/evidents [post]
func (h *Handler) CreateEvidentHandler(ctx *gin.Context) error {
	var user common.Users

	// ==========================
	// AUTH
	// ==========================
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

	// ==========================
	// FORM DATA
	// ==========================
	payload := &model.SCreateEvidentRequest{
		DealerID:      res.UserID,
		Category:      ctx.PostForm("category"),
		CatatanTemuan: ctx.PostForm("catatan_temuan"),
	}
	// Debug
	log.Printf("DealerID: %s", payload.DealerID)
	log.Printf("Category: %s", payload.Category)
	log.Printf("CatatanTemuan: %s", payload.CatatanTemuan)

	// ==========================
	// GET FILES
	// ==========================
	form, err := ctx.MultipartForm()
	if err != nil {
		return utils.ResponseError("99", err.Error(), ctx)
	}

	files := form.File["upload_temuan"]

	// ==========================
	// START TRANSACTION
	// ==========================
	tx := config.DB.WithContext(ctx.Request.Context()).Begin()

	if tx.Error != nil {
		return utils.ResponseError("99", tx.Error.Error(), ctx)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// ==========================
	// CREATE EVIDENT
	// ==========================
	newEvident := new(evident.Evident)

	payload.MapToSchema(newEvident, user.Username)

	if err := tx.Create(&newEvident).Error; err != nil {
		tx.Rollback()
		return utils.ResponseError("99", err.Error(), ctx)
	}

	// ==========================
	// LEADER BOARD
	// ==========================
	leaderBoardPoint := leaderboard.LeaderBoard{
		UserID:      user.ID,
		EvidentID:   &newEvident.ID,
		Point:       leaderboard.GetPointByCategory(newEvident.Category),
		Description: "Create Evident",
	}
	if err := tx.Create(&leaderBoardPoint).Error; err != nil {
		tx.Rollback()
		return utils.ResponseError("99", err.Error(), ctx)
	}
	// ==========================
	// SAVE FILES
	// ==========================
	for _, file := range files {

		filePath := utils.GenerateFilePathAndSave(
			file,
			"evident_before",
			ctx,
		)

		newPhoto := evidentPhoto.EvidentPhoto{
			EvidentID: &newEvident.ID,
			PhotoType: evidentPhoto.PHOTO_BEFORE,
			FilePath:  filePath,
		}

		if err := tx.Create(&newPhoto).Error; err != nil {
			tx.Rollback()

			_ = os.Remove(filePath)

			return utils.ResponseError("99", err.Error(), ctx)
		}
	}

	// ==========================
	// COMMIT
	// ==========================
	if err := tx.Commit().Error; err != nil {
		return utils.ResponseError("99", err.Error(), ctx)
	}

	if err := config.DB.Preload("Dealer").First(&newEvident, "id = ?", newEvident.ID).Error; err != nil {
		return utils.ResponseError("99", err.Error(), ctx)
	}

	return utils.ResponseSuccess("00", newEvident, ctx)
}

// @Tags Evident
// @Summary Revisi Evident
// @Description Revisi a evident entry
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Evident ID"
// @Param upload_temuan_after formData file true "Foto hasil revisi"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/evidents/{id} [patch]
func (h *Handler) RevisiEvidentHandler(ctx *gin.Context) error {
	var user common.Users

	id := ctx.Param("id")
	tokenString := ctx.GetHeader("Authorization")

	res, err := h.token.VerifyToken(tokenString)
	if err != nil {
		return utils.ResponseError("01", err.Error(), ctx)
	}

	if err := config.DB.
		Preload("Role").
		First(&user, "id = ?", res.UserID).
		Error; err != nil {
		return utils.ResponseError("99", "Unauthorized", ctx)
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		return utils.ResponseError("99", err.Error(), ctx)
	}

	files := form.File["upload_temuan_after"]

	if len(files) == 0 {
		return utils.ResponseError("99", "upload_temuan_after is required", ctx)
	}

	tx := config.DB.WithContext(ctx.Request.Context()).Begin()
	if tx.Error != nil {
		return utils.ResponseError("99", tx.Error.Error(), ctx)
	}

	var uploadedFiles []string

	defer func() {
		if r := recover(); r != nil {

			tx.Rollback()

			for _, file := range uploadedFiles {
				_ = os.Remove(file)
			}

			panic(r)
		}
	}()

	evidentID := uuid.MustParse(id)

	// ==========================
	// SAVE FILES
	// ==========================
	for _, file := range files {

		filePath := utils.GenerateFilePathAndSave(
			file,
			"evident_after",
			ctx,
		)

		if filePath == "" {
			tx.Rollback()

			for _, f := range uploadedFiles {
				_ = os.Remove(f)
			}

			return utils.ResponseError("99", "failed save file", ctx)
		}

		uploadedFiles = append(uploadedFiles, filePath)

		newPhoto := evidentPhoto.EvidentPhoto{
			EvidentID: &evidentID,
			PhotoType: evidentPhoto.PHOTO_AFTER,
			FilePath:  filePath,
		}

		if err := tx.Create(&newPhoto).Error; err != nil {

			tx.Rollback()

			for _, f := range uploadedFiles {
				_ = os.Remove(f)
			}

			return utils.ResponseError("99", err.Error(), ctx)
		}
	}

	// ==========================
	// UPDATE STATUS
	// ==========================
	if err := tx.
		Model(&evident.Evident{}).
		Where("id = ?", id).
		Update("status", evident.StatusCompleted).
		Error; err != nil {

		tx.Rollback()

		for _, f := range uploadedFiles {
			_ = os.Remove(f)
		}

		return utils.ResponseError("99", err.Error(), ctx)
	}

	if err := tx.Commit().Error; err != nil {

		for _, f := range uploadedFiles {
			_ = os.Remove(f)
		}

		return utils.ResponseError("99", err.Error(), ctx)
	}

	return utils.ResponseSuccess("00", gin.H{
		"uploaded_files": len(uploadedFiles),
	}, ctx)
}

// @Tags Evident
// @Summary Get All Evident.
// @Description Retrieve all evident entries.
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /v1/evidents [get]
// @securitydefinitions.apikey token
// @in header
// @name Authorization
// @Security token
func (h *Handler) GetAllEvident(ctx *gin.Context) error {
	var user common.Users
	tokenString := ctx.GetHeader("Authorization")
	var evidents []evident.Evident

	res, err := h.token.VerifyToken(tokenString)
	if err != nil {
		return utils.ResponseError("01", err.Error(), ctx)
	}
	if err := config.DB.Preload("Role").First(&user, "id = ?", res.UserID).Error; err != nil {
		return utils.ResponseError("99", "Unauthorized", ctx)
	}
	log.Printf("User: %+v", user)
	if user.Role.Name == "AUDIT" {
		if err := config.DB.Preload("Dealer").Find(&evidents).Error; err != nil {
			return utils.ResponseError("99", err.Error(), ctx)
		}
		return utils.ResponseSuccess("00", evidents, ctx)
	}

	if err := config.DB.Preload("Dealer").Where("dealer_id = ?", res.UserID).Find(&evidents).Error; err != nil {
		return utils.ResponseError("99", err.Error(), ctx)
	}

	return utils.ResponseSuccess("00", evidents, ctx)
}

// @Tags Evident
// @Summary Get By Id Evident.
// @Description Retrieve one evident entries.
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /v1/evidents/{id} [get]
// @securitydefinitions.apikey token
// @in header
// @name Authorization
// @Security token
func (h *Handler) GetEvidentById(ctx *gin.Context) error {
	var user common.Users
	id := ctx.Param("id")
	tokenString := ctx.GetHeader("Authorization")
	var evident evident.Evident

	res, err := h.token.VerifyToken(tokenString)
	if err != nil {
		return utils.ResponseError("01", err.Error(), ctx)
	}
	if err := config.DB.Preload("Role").First(&user, "id = ?", res.UserID).Error; err != nil {
		return utils.ResponseError("99", "Unauthorized", ctx)
	}

	if user.Role.Name == "AUDIT" {
		if err := config.DB.Preload("Dealer").Preload("Photos").First(&evident, "id = ?", id).Error; err != nil {
			return utils.ResponseError("99", err.Error(), ctx)
		}
		return utils.ResponseSuccess("00", evident, ctx)
	} else {

		if err := config.DB.Preload("Dealer").Preload("Photos").First(&evident, "id = ?", id).Error; err != nil {
			return utils.ResponseError("99", err.Error(), ctx)
		}

		return utils.ResponseSuccess("00", evident, ctx)
	}
}
