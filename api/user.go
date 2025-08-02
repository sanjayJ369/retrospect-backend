package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/sanjayj369/retrospect-backend/db/sqlc"
	"github.com/sanjayj369/retrospect-backend/token"
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

	err = authorizeUser(ctx, user.ID.Bytes)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
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
	SessionId             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
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

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(userId, s.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(userId, s.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           pgtype.UUID{Bytes: refreshPayload.ID, Valid: true},
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    pgtype.Timestamp{Time: refreshPayload.ExpiredAt, Valid: true},
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	res := UserLoginResponse{
		SessionId:             session.ID.Bytes,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, res)

}

func authorizeUser(ctx *gin.Context, userId uuid.UUID) error {
	tkn := ctx.MustGet(authorizationPayloadKey)
	payload, ok := tkn.(*token.Payload)
	if !ok {
		err := fmt.Errorf("invalid token")
		return err
	}
	if payload.UserId != userId {
		err := fmt.Errorf("account does not belong to the authorized user")
		return err
	}
	return nil
}
