package interfaces

import (
	"context"

	"github.com/SethukumarJ/CashierX-Auth-Service/pkg/domain"
)

type UserRepository interface {
	FindAll(ctx context.Context) ([]domain.Users, error)
	FindByID(ctx context.Context, id uint) (domain.Users, error)
	FindByName(ctx context.Context, email string) (domain.Users, error)
	Save(ctx context.Context, user domain.Users) (domain.Users, error)
	Delete(ctx context.Context, user domain.Users) error
	FindPassword(ctx context.Context, email string) (string, error)
	
}
