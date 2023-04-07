package interfaces

import (
	"context"

	domain "github.com/SethukumarJ/CashierX-Auth-Service/pkg/domain"
)

type UserUseCase interface {
	FindAll(ctx context.Context) ([]domain.Users, error)
	FindByID(ctx context.Context, id uint) (domain.Users, error)
	FindByName(ctx context.Context, email string) (domain.Users, error)
	Register(ctx context.Context, user domain.Users) (domain.Users, error)
	Delete(ctx context.Context, user domain.Users) error
	VerifyUser(ctx context.Context, email string, password string) error
}
