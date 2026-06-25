package middleware

import (
	"evidence-hub-be/src/core/pkg/token"
)

type Middleware struct {
	token *token.Token
}

func New(token *token.Token) *Middleware {
	return &Middleware{
		token: token,
	}
}
