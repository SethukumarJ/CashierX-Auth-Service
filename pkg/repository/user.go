package repository

import (
	"context"
	"fmt"

	domain "github.com/SethukumarJ/CashierX-Auth-Service/pkg/domain"
	interfaces "github.com/SethukumarJ/CashierX-Auth-Service/pkg/repository/interface"
	"gorm.io/gorm"
)

type userDatabase struct {
	DB *gorm.DB
}

// FindPassword implements interfaces.UserRepository
func (c *userDatabase) FindPassword(ctx context.Context, email string) (string, error) {
	var user domain.Users
	err := c.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return "", err
	}

	return user.Password, nil
}

// FindByName implements interfaces.UserRepository
func (c *userDatabase) FindByName(ctx context.Context, email string) (domain.Users, error) {
	var user domain.Users
	err := c.DB.Where("email = ?", email).First(&user).Error
	fmt.Println("err from repos find name", err)
	return user, err
}

func NewUserRepository(DB *gorm.DB) interfaces.UserRepository {
	return &userDatabase{DB}
}

func (c *userDatabase) FindAll(ctx context.Context) ([]domain.Users, error) {
	var users []domain.Users
	err := c.DB.Find(&users).Error

	return users, err
}

func (c *userDatabase) FindByID(ctx context.Context, id uint) (domain.Users, error) {
	var user domain.Users
	err := c.DB.First(&user, id).Error

	return user, err
}

func (c *userDatabase) Save(ctx context.Context, user domain.Users) (domain.Users, error) {
	err := c.DB.Save(&user).Error

	return user, err
}

func (c *userDatabase) Delete(ctx context.Context, user domain.Users) error {
	err := c.DB.Delete(&user).Error

	return err
}
