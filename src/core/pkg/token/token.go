package token

import (
	"evidence-hub-be/src/core/config/envs"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type Token struct {
	cfg *envs.Config
}

func New(cfg *envs.Config) *Token {
	return &Token{
		cfg: cfg,
	}
}

func (t *Token) GenerateToken(payload *Payload) (string, error) {
	tokenByte := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)
	jwt_exp, _ := time.ParseDuration(t.cfg.AppConfig.JwtExpiredIn)
	claims["id_user"] = payload.UserID
	claims["id_role"] = payload.Role.ID
	claims["role"] = payload.Role.Name
	claims["role_type"] = payload.Role.Type
	claims["exp"] = now.Add(jwt_exp).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()
	if payload.PrivilageAdmin != "" {
		claims["privilage"] = payload.PrivilageAdmin
	}

	tokenString, err := tokenByte.SignedString([]byte(t.cfg.AppConfig.JwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t *Token) VerifyToken(tokenString string) (*Response, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.cfg.AppConfig.JwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	userID, userIDExists := claims["id_user"].(string)

	if !userIDExists {
		return nil, fmt.Errorf("user not found")
	}

	roleID, roleIDExists := claims["id_role"].(string)

	if !roleIDExists {
		return nil, fmt.Errorf("role not found")
	}

	role, roleExists := claims["role"].(string)

	if !roleExists {
		return nil, fmt.Errorf("role not found")
	}

	roleType, roleTypeExists := claims["role_type"].(string)

	if !roleTypeExists {
		return nil, fmt.Errorf("role type not found")
	}

	privilageAdmin, _ := claims["privilage"].(string)

	return &Response{
		UserID:    userID,
		RoleID:    roleID,
		Role:      role,
		RoleType:  roleType,
		Privilage: privilageAdmin,
	}, nil
}
