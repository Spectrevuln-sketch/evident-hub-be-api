package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseErrorType struct {
	RespCode string `json:"resp_code"`
	Message  string `json:"message"`
}

type ResponseSuccessType struct {
	RespCode    string      `json:"resp_code"`
	Data        interface{} `json:"data"`
	ResponsCode string      `json:"respons_code,omitempty"`
}

func ResponseError(code string, errMsg string, ctx *gin.Context) error {
	switch code {
	case "01":
		ctx.JSON(http.StatusUnauthorized, ResponseErrorType{
			RespCode: code,
			Message:  errMsg,
		})
	case "05":
		ctx.JSON(http.StatusForbidden, ResponseErrorType{
			RespCode: code,
			Message:  errMsg,
		})
	default:
		ctx.JSON(http.StatusBadRequest, ResponseErrorType{
			RespCode: code,
			Message:  errMsg,
		})
	}

	return nil
}

func ResponseSuccess(code string, data interface{}, ctx *gin.Context) error {
	ctx.JSON(http.StatusOK, ResponseSuccessType{
		RespCode:    code,
		ResponsCode: "OK",
		Data:        data,
	})

	return nil
}

func Wrap(h func(*gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = h(c)
	}
}
