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
		Email:    req.Email,
		Name:     req.Name,
		Password: hashedPassword,
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

type getUsersRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type createUserResponse struct {
	ID        pgtype.UUID      `json:"id"`
	Email     string           `json:"email"`
	Name      string           `json:"name"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	Timezone  string           `json:"timezone"`
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

	res := createUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Timezone:  user.Timezone,
	}

	ctx.JSON(http.StatusOK, res)
}
