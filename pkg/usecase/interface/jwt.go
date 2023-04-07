package interfaces

import (
	"github.com/golang-jwt/jwt"
	domain "github.com/SethukumarJ/CashierX-Auth-Service/pkg/domain"
)

type JWTUsecase interface {
	GenerateAccessToken(userid uint, userName string) (string, error)
	VerifyToken(token string) (bool, *domain.SignedDetails)
	GetTokenFromString(signedToken string, claims *domain.SignedDetails) (*jwt.Token, error)
	GenerateRefreshToken(userid uint, userName string) (string, error)

}
