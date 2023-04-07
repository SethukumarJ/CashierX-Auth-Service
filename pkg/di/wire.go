//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	http "github.com/SethukumarJ/CashierX-Auth-Service/pkg/api"
	handler "github.com/SethukumarJ/CashierX-Auth-Service/pkg/api/handler"
	config "github.com/SethukumarJ/CashierX-Auth-Service/pkg/config"
	db "github.com/SethukumarJ/CashierX-Auth-Service/pkg/db"
	repository "github.com/SethukumarJ/CashierX-Auth-Service/pkg/repository"
	usecase "github.com/SethukumarJ/CashierX-Auth-Service/pkg/usecase"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	wire.Build(db.ConnectDatabase, 
			repository.NewUserRepository, 
			usecase.NewUserUseCase, 
			handler.NewUserHandler, 
			http.NewServerHTTP)

	return &http.ServerHTTP{}, nil
}
