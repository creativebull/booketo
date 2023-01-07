package controllers

import (
	"net/http"

	usersDomain "github.com/Ayobami-00/booketo-mvc-go-postgres-gin/src/domain/users"
	services "github.com/Ayobami-00/booketo-mvc-go-postgres-gin/src/services/users"
	apiResponse "github.com/Ayobami-00/booketo-mvc-go-postgres-gin/src/utils/api_response"
	"github.com/gin-gonic/gin"
)

const (
	userCreateSuccessMessage = "User created successfully"
)

func CreateUser(ctx *gin.Context) {
	var request usersDomain.CreateUserRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, apiResponse.NewErrorApiResponse(err.Error()))
		return
	}

	result, err := services.UsersService.CreateUser(ctx, request)

	if err != nil {
		ctx.JSON(err.Status(), apiResponse.NewErrorApiResponse(err.Message()))
		return
	}

	ctx.JSON(http.StatusCreated, apiResponse.NewSuccessApiResponseWithData(userCreateSuccessMessage, result))
}
