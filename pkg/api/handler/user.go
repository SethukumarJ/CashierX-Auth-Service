package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	domain "github.com/SethukumarJ/CashierX-Auth-Service/pkg/domain"
	"github.com/SethukumarJ/CashierX-Auth-Service/pkg/pb"
	services "github.com/SethukumarJ/CashierX-Auth-Service/pkg/usecase/interface"
)

type UserHandler struct {
	userUseCase services.UserUseCase
	jwtUsecase  services.JWTUsecase
}

type Response struct {
	Id       int64  `copier:"must"`
	Email    string `copier:"must"`
	Password string `copier:"must"`
}

func NewUserHandler(usecase services.UserUseCase, jwtusecase services.JWTUsecase) *UserHandler {
	return &UserHandler{
		userUseCase: usecase,
		jwtUsecase:  jwtusecase,
	}
}

func (cr *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user := domain.Users{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		UserName:  req.UserName,
	}

	// check if email already exists
	user1, err := cr.userUseCase.FindByName(ctx, user.Email)
	if err == nil {
		return &pb.RegisterResponse{
			Status: http.StatusUnprocessableEntity,
			Id:     user1.Id,
			Error:  fmt.Sprint(errors.New("email already exists")),
		}, nil
	}

	// register the user
	user, err = cr.userUseCase.Register(ctx, user)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return &pb.RegisterResponse{
				Status: http.StatusConflict,
				Error:  fmt.Sprint(errors.New("username already exists")),
			}, nil
		}
		return &pb.RegisterResponse{
			Status: http.StatusUnprocessableEntity,
			Error:  fmt.Sprint(errors.New("failed to register user")),
		}, nil
	}

	return &pb.RegisterResponse{
		Status: http.StatusCreated,
		Id:     user.Id,
	}, nil
}

func (cr *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	err := cr.userUseCase.VerifyUser(ctx, req.Email, req.Password)
	if err != nil {
		return &pb.LoginResponse{
			Status: http.StatusUnauthorized,
			Error:  fmt.Sprintf("failed to verify user: %s", err.Error()),
		}, nil
	}

	user, err := cr.userUseCase.FindByName(ctx, req.Email)
	if err != nil {
		return &pb.LoginResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("error while getting user from db: %s", err.Error()),
		}, nil
	}
	accesstoken, err := cr.jwtUsecase.GenerateAccessToken(uint(user.Id), user.Email)
	if err != nil {
		return &pb.LoginResponse{
			Status: http.StatusUnauthorized,
			Error:  fmt.Sprint(errors.New("failed to generate access token")),
		}, errors.New(err.Error())
	}
	refreshtoken, err := cr.jwtUsecase.GenerateRefreshToken(uint(user.Id), user.Email)
	if err != nil {
		return &pb.LoginResponse{
			Status: http.StatusUnauthorized,
			Error:  fmt.Sprint(errors.New("failed to generate refresh token")),
		}, errors.New(err.Error())
	}
	return &pb.LoginResponse{
		Status:        http.StatusOK,
		AccessToken:   accesstoken,
		RefresshToken: refreshtoken,
	}, nil
}

func (cr *UserHandler) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	ok, claims := cr.jwtUsecase.VerifyToken(req.Token)
	fmt.Println("claims", claims)
	if !ok {
		return &pb.ValidateResponse{
			Status: http.StatusUnauthorized,
			Error:  fmt.Sprint(errors.New("token verification failed")),
		}, nil
	}

	user, err := cr.userUseCase.FindByName(ctx, claims.UserName)

	if err != nil {
		return &pb.ValidateResponse{
			Status: http.StatusUnauthorized,
			Error:  fmt.Sprint(errors.New("user not found with token credentials")),
		}, errors.New(err.Error())
	}

	return &pb.ValidateResponse{
		Status: http.StatusOK,
		UserId: user.Id,
		Source: fmt.Sprint(claims.Source),
	}, nil
}

func (cr *UserHandler) TokenRefresh(ctx context.Context, req *pb.TokenRefreshRequest) (*pb.TokenRefreshResponse, error) {

	ok, claims := cr.jwtUsecase.VerifyToken(req.Token)
	if !ok {
		return &pb.TokenRefreshResponse{
			Status: http.StatusUnauthorized,
			Error:  fmt.Sprint(errors.New("token verification failed")),
		}, nil
	}

	fmt.Println("//////////////////////////////////", claims.UserName)
	accesstoken, err := cr.jwtUsecase.GenerateAccessToken(claims.UserId, claims.UserName)

	if err != nil {
		return &pb.TokenRefreshResponse{
			Status: http.StatusUnauthorized,
			Error:  fmt.Sprint(errors.New("unable to generate access token")),
		}, errors.New(err.Error())
	}
	return &pb.TokenRefreshResponse{
		Status: http.StatusOK,
		Token:  accesstoken,
	}, nil

}

func (cr *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	// Check if the ID is not empty or invalid
	if req.Id == 0 {
		return &pb.DeleteUserResponse{
			Status: http.StatusBadRequest,
			Error:  "Invalid ID",
		}, nil
	}

	var user domain.Users

	// Check if the record exists in the database
	user, err := cr.userUseCase.FindByID(ctx, uint(req.Id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.DeleteUserResponse{
				Status: http.StatusNotFound,
				Error:  "Record not found",
			}, nil
		} else {
			return &pb.DeleteUserResponse{
				Status: http.StatusInternalServerError,
				Error:  fmt.Sprint(errors.New("unable to fetch user")),
			}, nil
		}
	}

	// Delete the record from the database
	err = cr.userUseCase.Delete(ctx, int64(req.Id))
	if err != nil {
		return &pb.DeleteUserResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprint(errors.New("error in deleting data")),
		}, nil
	}

	return &pb.DeleteUserResponse{
		Status: http.StatusOK,
		Id:     user.Id,
	}, nil
}

func (cr *UserHandler) FindUser(ctx context.Context, req *pb.FindUserRequest) (*pb.FindUserResponse, error) {
	// Check if the ID is not empty or invalid
	if req.Id == 0 {
		return &pb.FindUserResponse{
			Status: http.StatusBadRequest,
			Error:  "Invalid ID",
		}, nil
	}

	var user domain.Users

	// Check if the record exists in the database
	user, err := cr.userUseCase.FindByID(ctx, uint(req.Id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.FindUserResponse{
				Status: http.StatusNotFound,
				Error:  "Record not found",
			}, nil
		} else {
			return &pb.FindUserResponse{
				Status: http.StatusInternalServerError,
				Error:  fmt.Sprint(errors.New("unable to fetch user")),
			}, nil
		}
	}
	data := &pb.FindUser{
		Id:        user.Id,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		UserName:  user.UserName,
	}
	return &pb.FindUserResponse{
		Status: http.StatusOK,
		Data:   data,
	}, nil
}

func (cr *UserHandler) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	// Get all users from the use case
	users, err := cr.userUseCase.FindAll(ctx)
	if err != nil {
		return &pb.GetUsersResponse{
			Status: http.StatusUnprocessableEntity,
			Error:  fmt.Sprint(errors.New("unable to fetch data")),
		}, errors.New(err.Error())
	}

	// Convert domain.Users to pb.User
	var pbUsers []*pb.User
	for _, user := range users {
		pbUser := &pb.User{
			Id:        user.Id,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			UserName:  user.UserName,
		}
		pbUsers = append(pbUsers, pbUser)
	}

	// Return the response
	return &pb.GetUsersResponse{
		Status: http.StatusOK,
		User:   pbUsers,
	}, nil
}

func (cr *UserHandler) FindUsers(c *gin.Context) {
	paramsId := c.Param("id")
	id, err := strconv.Atoi(paramsId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "cannot parse id",
		})
		return
	}

	user, err := cr.userUseCase.FindByID(c.Request.Context(), uint(id))

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		response := Response{}
		copier.Copy(&response, &user)

		c.JSON(http.StatusOK, response)
	}
}
