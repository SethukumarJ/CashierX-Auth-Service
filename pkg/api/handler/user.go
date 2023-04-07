package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"

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

func NewUserHandler(usecase services.UserUseCase,jwtusecase services.JWTUsecase) *UserHandler {
	return &UserHandler{
		userUseCase: usecase,
		jwtUsecase: jwtusecase,
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

	user1, err := cr.userUseCase.FindByName(ctx, user.Email)
	if err == nil {
		fmt.Println(errors.New("email already exist"))
		return &pb.RegisterResponse{
			Status: http.StatusUnprocessableEntity,
			Id:     user1.Id,
			Error: fmt.Sprint(errors.New("email already exist")),
		}, err
	}
	user, err = cr.userUseCase.Register(ctx, user)
	fmt.Println(user)
	if err != nil {
		fmt.Println(errors.New("email already exist//////"))
		return nil, err
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
	accesstoken, err := cr.jwtUsecase.GenerateAccessToken(uint(user.Id),user.Email)
	if err != nil {
		return &pb.LoginResponse{
			Status: http.StatusUnauthorized,
			Error: fmt.Sprint(errors.New("failed to generate access token")),
		}, errors.New(err.Error())
	}
	refreshtoken, err := cr.jwtUsecase.GenerateRefreshToken(uint(user.Id),user.Email)
	if err != nil {
		return &pb.LoginResponse{
			Status: http.StatusUnauthorized,
			Error: fmt.Sprint(errors.New("failed to generate refresh token")),
		}, errors.New(err.Error())
	}
	return &pb.LoginResponse{
		Status: http.StatusOK,
		AccessToken: accesstoken,
		RefresshToken: refreshtoken,
	}, nil
}

func (cr *UserHandler) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	// claims, err := s.Jwt.ValidateToken(req.Token)

	// if err != nil {
	// 	return &pb.ValidateResponse{
	// 		Status: http.StatusBadRequest,
	// 		Error:  err.Error(),
	// 	}, nil
	// }

	var user domain.Users

	// if result := s.H.DB.Where(&models.User{Email: claims.Email}).First(&user); result.Error != nil {
	// 	return &pb.ValidateResponse{
	// 		Status: http.StatusNotFound,
	// 		Error:  "User not found",
	// 	}, nil
	// }

	return &pb.ValidateResponse{
		Status: http.StatusOK,
		UserId: user.Id,
	}, nil
}

func (cr *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	// // Check if the ID is not empty or invalid
	// if req.Id == 0 {
	// 	return &pb.DeleteResponse{
	// 		Status: http.StatusBadRequest,
	// 		Error:  "Invalid ID",
	// 	}, nil
	// }

	var user domain.Users

	// // Check if the record exists in the database
	// result := s.H.DB.First(&user, "id = ?", req.Id)
	// if result.Error != nil {
	// 	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
	// 		return &pb.DeleteResponse{
	// 			Status: http.StatusNotFound,
	// 			Error:  "Record not found",
	// 		}, nil
	// 	} else {
	// 		return &pb.DeleteResponse{
	// 			Status: http.StatusInternalServerError,
	// 			Error:  result.Error.Error(),
	// 		}, nil
	// 	}
	// }

	// // Delete the record from the database
	// result = s.H.DB.Delete(&user, req.Id)
	// if result.Error != nil {
	// 	return &pb.DeleteResponse{
	// 		Status: http.StatusInternalServerError,
	// 		Error:  result.Error.Error(),
	// 	}, nil
	// }

	return &pb.DeleteUserResponse{
		Status: http.StatusOK,
		Id:     user.Id,
	}, nil
}

// FindAll godoc
// @summary Get all users
// @description Get all users
// @tags users
// @security ApiKeyAuth
// @id FindAll
// @produce json
// @Router /api/users [get]
// @response 200 {object} []Response "OK"
func (cr *UserHandler) FindAll(c *gin.Context) {
	users, err := cr.userUseCase.FindAll(c.Request.Context())

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		response := []Response{}
		copier.Copy(&response, &users)

		c.JSON(http.StatusOK, response)
	}
}

func (cr *UserHandler) FindUser(ctx context.Context, req *pb.FindUserRequest) (*pb.FindUserResponse, error) {
	// // Check if the ID is not empty or invalid
	// if req.Id == 0 {
	// 	return &pb.DeleteResponse{
	// 		Status: http.StatusBadRequest,
	// 		Error:  "Invalid ID",
	// 	}, nil
	// }

	// var user domain.Users

	// // Check if the record exists in the database
	// result := s.H.DB.First(&user, "id = ?", req.Id)
	// if result.Error != nil {
	// 	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
	// 		return &pb.DeleteResponse{
	// 			Status: http.StatusNotFound,
	// 			Error:  "Record not found",
	// 		}, nil
	// 	} else {
	// 		return &pb.DeleteResponse{
	// 			Status: http.StatusInternalServerError,
	// 			Error:  result.Error.Error(),
	// 		}, nil
	// 	}
	// }

	// // Delete the record from the database
	// result = s.H.DB.Delete(&user, req.Id)
	// if result.Error != nil {
	// 	return &pb.DeleteResponse{
	// 		Status: http.StatusInternalServerError,
	// 		Error:  result.Error.Error(),
	// 	}, nil
	// }

	return &pb.FindUserResponse{
		Status: http.StatusOK,
	}, nil
}

func (cr *UserHandler) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	// // Check if the ID is not empty or invalid
	// if req.Id == 0 {
	// 	return &pb.DeleteResponse{
	// 		Status: http.StatusBadRequest,
	// 		Error:  "Invalid ID",
	// 	}, nil
	// }

	// var user domain.Users

	// // Check if the record exists in the database
	// result := s.H.DB.First(&user, "id = ?", req.Id)
	// if result.Error != nil {
	// 	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
	// 		return &pb.DeleteResponse{
	// 			Status: http.StatusNotFound,
	// 			Error:  "Record not found",
	// 		}, nil
	// 	} else {
	// 		return &pb.DeleteResponse{
	// 			Status: http.StatusInternalServerError,
	// 			Error:  result.Error.Error(),
	// 		}, nil
	// 	}
	// }

	// // Delete the record from the database
	// result = s.H.DB.Delete(&user, req.Id)
	// if result.Error != nil {
	// 	return &pb.DeleteResponse{
	// 		Status: http.StatusInternalServerError,
	// 		Error:  result.Error.Error(),
	// 	}, nil
	// }

	return &pb.GetUsersResponse{}, nil
}

func (cr *UserHandler) TokenRefresh(ctx context.Context, req *pb.TokenRefreshRequest) (*pb.TokenRefreshResponse, error) {
	// // Check if the ID is not empty or invalid
	// if req.Id == 0 {
	// 	return &pb.DeleteResponse{
	// 		Status: http.StatusBadRequest,
	// 		Error:  "Invalid ID",
	// 	}, nil
	// }

	// var user domain.Users

	// // Check if the record exists in the database
	// result := s.H.DB.First(&user, "id = ?", req.Id)
	// if result.Error != nil {
	// 	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
	// 		return &pb.DeleteResponse{
	// 			Status: http.StatusNotFound,
	// 			Error:  "Record not found",
	// 		}, nil
	// 	} else {
	// 		return &pb.DeleteResponse{
	// 			Status: http.StatusInternalServerError,
	// 			Error:  result.Error.Error(),
	// 		}, nil
	// 	}
	// }

	// // Delete the record from the database
	// result = s.H.DB.Delete(&user, req.Id)
	// if result.Error != nil {
	// 	return &pb.DeleteResponse{
	// 		Status: http.StatusInternalServerError,
	// 		Error:  result.Error.Error(),
	// 	}, nil
	// }

	return &pb.TokenRefreshResponse{}, nil
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
