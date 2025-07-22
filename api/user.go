package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/sanjayj369/retrospect-backend/db/sqlc"
	"github.com/sanjayj369/retrospect-backend/util"
)

type createUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashedPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Email:          req.Email,
		Name:           req.Name,
		HashedPassword: hashedPassword,
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, newUserResponse(user))
}

type getUsersRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type userResponse struct {
	ID        pgtype.UUID      `json:"id"`
	Email     string           `json:"email"`
	Name      string           `json:"name"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	Timezone  string           `json:"timezone"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Timezone:  user.Timezone,
	}
}
func (s *Server) getUser(ctx *gin.Context) {
	var req getUsersRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	parsedUUID, _ := uuid.Parse(req.ID)

	var userUUIDBytes [16]byte = parsedUUID

	user, err := s.store.GetUser(ctx, pgtype.UUID{Bytes: userUUIDBytes, Valid: true})
	if err != nil {
		fmt.Printf("%v", err)
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := newUserResponse(user)

	ctx.JSON(http.StatusOK, res)
}

type UserLoginRequest struct {
	Name     string `json:"name" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserLoginResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (s *Server) LoginUser(ctx *gin.Context) {
	var req UserLoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	user, err := s.store.GetUserByName(ctx, req.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, nil)
		return
	}

	userId, err := uuid.FromBytes(user.ID.Bytes[:])
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	tkn, err := s.tokenMaker.CreateToken(userId, s.config.Duration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	res := UserLoginResponse{
		AccessToken: tkn,
		User:        newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, res)

}
