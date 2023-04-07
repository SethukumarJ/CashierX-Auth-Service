package http

import (
	"fmt"
	"log"
	"net"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	pb "github.com/SethukumarJ/CashierX-Auth-Service/pkg/pb"
	"google.golang.org/grpc"

	_ "github.com/SethukumarJ/CashierX-Auth-Service/cmd/api/docs"
	handler "github.com/SethukumarJ/CashierX-Auth-Service/pkg/api/handler"
	middleware "github.com/SethukumarJ/CashierX-Auth-Service/pkg/api/middleware"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func StartGRPCServer(userHandler *handler.UserHandler, grpcPort string) {
	lis, err := net.Listen("tcp", ":"+grpcPort)
	fmt.Println("grpcPort/////", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, userHandler)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func NewServerHTTP(userHandler *handler.UserHandler) *ServerHTTP {
	engine := gin.New()
	go StartGRPCServer(userHandler, "50056")
	// Use logger from Gin
	engine.Use(gin.Logger())

	// Swagger docs
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Request JWT
	engine.POST("/login", middleware.LoginHandler)

	// Auth middleware
	// api := engine.Group("/api", middleware.AuthorizationMiddleware)

	// engine.GET("/api/users", userHandler.FindAll)
	// // api.GET("users/:id", userHandler.FindByID)
	// engine.POST("/api/users", userHandler.Save)
	// api.DELETE("users/:id", userHandler.Delete)

	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Start() {
	sh.engine.Run(":3002")
}
