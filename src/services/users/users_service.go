package services

import (
	"context"
	"database/sql"

	app "github.com/Ayobami-00/booketo-mvc-go-postgres-gin/src/app"
	store "github.com/Ayobami-00/booketo-mvc-go-postgres-gin/src/db"
	db "github.com/Ayobami-00/booketo-mvc-go-postgres-gin/src/db/sqlc"
	"github.com/Ayobami-00/booketo-mvc-go-postgres-gin/src/domain/users"
	"github.com/Ayobami-00/booketo-mvc-go-postgres-gin/src/utils"
	errors "github.com/Ayobami-00/booketo-mvc-go-postgres-gin/src/utils/errors"
	"github.com/lib/pq"
)

type usersService struct {
	DbStore db.Store
}

type UserSessionData struct {
	UserAgent string
	ClientIP  string
}

type usersServiceInterface interface {
	CreateUser(ctx context.Context, request users.CreateUserRequest) (*users.UserResponse, errors.ApiError)
	LoginUser(ctx context.Context, request users.LoginUserRequest, sessionData UserSessionData) (*users.LoginUserResponse, errors.ApiError)
}

var (
	UsersService usersServiceInterface
)

const (

	// USER ERRORS

	userCreateGenericError = "Something went wrong when creating user"
	userLoginGenericError  = "Something went wrong when login in user"
)

func init() {
	UsersService = &usersService{
		DbStore: store.DbStore,
	}
}

func (s *usersService) CreateUser(ctx context.Context, request users.CreateUserRequest) (*users.UserResponse, errors.ApiError) {

	hashedPassword, err := utils.HashPassword(request.Password)

	if err != nil {
		return nil, errors.NewInternalServerError(userCreateGenericError)
	}

	arg := db.CreateUserParams{
		Email:          request.Email,
		HashedPassword: hashedPassword,
	}

	createdUser, err := s.DbStore.CreateUser(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, errors.NewForbiddenError(err.Error())
			}
		}

		return nil, errors.NewInternalServerError(userCreateGenericError)
	}

	return users.NewCreateUserResponse(createdUser), nil

}

func (s *usersService) LoginUser(ctx context.Context, request users.LoginUserRequest, sessionData UserSessionData) (*users.LoginUserResponse, errors.ApiError) {

	user, err := s.DbStore.GetUser(ctx, request.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError(err.Error())
		}
		errors.NewInternalServerError(userLoginGenericError)
	}

	err = utils.CheckPassword(request.Password, user.HashedPassword)
	if err != nil {
		return nil, errors.NewNotFoundError(err.Error())
	}

	accessToken, accessPayload, err := app.TokenMaker.CreateToken(
		user.ID,
		app.Config.AccessTokenDuration,
	)
	if err != nil {
		return nil, errors.NewInternalServerError(userLoginGenericError)
	}

	refreshToken, refreshPayload, err := app.TokenMaker.CreateToken(
		user.ID,
		app.Config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, errors.NewInternalServerError(userLoginGenericError)
	}

	session, err := s.DbStore.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    sessionData.UserAgent,
		ClientIp:     sessionData.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, errors.NewInternalServerError(userLoginGenericError)
	}

	rsp := &users.LoginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  users.NewCreateUserResponse(user),
	}

	return rsp, nil

}
