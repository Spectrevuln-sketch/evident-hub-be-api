package users

import (
	"evidence-hub-be/src/core/config/envs"
	"evidence-hub-be/src/core/pkg/token"
)

type Handler struct {
	cfg   *envs.Config
	token *token.Token
}

func NewHandler(cfg *envs.Config, token *token.Token) *Handler {
	return &Handler{
		cfg:   cfg,
		token: token,
	}
}
