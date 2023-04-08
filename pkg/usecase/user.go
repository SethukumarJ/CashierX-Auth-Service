package usecase

import (
	"context"
	"errors"
	"fmt"
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
func (cr *userUseCase) Delete(ctx context.Context, id int64) error {
	err := cr.userRepo.Delete(ctx,id)
	if err != nil {
		return err
	}
	return nil
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

// / VerifyUser verifies the user credentials
func (c *userUseCase) VerifyUser(ctx context.Context, email string, password string) error {

	_, err := c.userRepo.FindByName(ctx, email)

	if err != nil {
		fmt.Println(errors.New("failed to login. check your email"))
		return errors.New("failed to login. check your email")
	}
	pswd, err := c.userRepo.FindPassword(ctx, email)
	if err != nil {
		return errors.New("failed to login. check your email or password")
	}
	isValidPassword := VerifyPassword(pswd, []byte(password))
	if !isValidPassword {
		return errors.New("failed to login. check your credential")
	}

	return nil
}

func VerifyPassword(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
