package usecase

import (
	"context"
	"log"

	domain "github.com/SethukumarJ/CashierX-Auth-Service/pkg/domain"
	interfaces "github.com/SethukumarJ/CashierX-Auth-Service/pkg/repository/interface"
	services "github.com/SethukumarJ/CashierX-Auth-Service/pkg/usecase/interface"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	userRepo interfaces.UserRepository
}

// FindByName implements interfaces.UserUseCase
func (c *userUseCase) FindByName(ctx context.Context, email string) (domain.Users, error) {
	user, err := c.userRepo.FindByName(ctx, email)
	return user, err
}

// Delete implements interfaces.UserUseCase
func (*userUseCase) Delete(ctx context.Context, user domain.Users) error {
	panic("unimplemented")
}

func NewUserUseCase(repo interfaces.UserRepository) services.UserUseCase {
	return &userUseCase{
		userRepo: repo,
	}
}

func (c *userUseCase) FindAll(ctx context.Context) ([]domain.Users, error) {
	users, err := c.userRepo.FindAll(ctx)
	return users, err
}

func (c *userUseCase) FindByID(ctx context.Context, id uint) (domain.Users, error) {
	user, err := c.userRepo.FindByID(ctx, id)
	return user, err
}

func (c *userUseCase) Register(ctx context.Context, user domain.Users) (domain.Users, error) {

	user.Password = HashPassword(user.Password)
	user, err := c.userRepo.Save(ctx, user)

	return user, err
}

// HashPassword hashes the password
func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	return string(hash)
}

// func (c *userUseCase) Delete(ctx context.Context, user domain.Users) error {
// 	err := c.userRepo.Delete(ctx, user)

// 	return err
// }
