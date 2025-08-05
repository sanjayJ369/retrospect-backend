package api

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/sanjayj369/retrospect-backend/db/sqlc"
	"github.com/sanjayj369/retrospect-backend/mail"
	"github.com/sanjayj369/retrospect-backend/token"
	"github.com/sanjayj369/retrospect-backend/util"
)

type resetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

func (s *Server) resetPassword(ctx *gin.Context) {
	var req resetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashedPassword(req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	payload, err := s.tokenMaker.VerifyToken(req.Token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if payload.Purpose != token.PurposeResetPassword {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid token purpose")))
		return
	}

	arg := db.UpdateUserHashedPasswordParams{
		ID:             pgtype.UUID{Bytes: payload.UserId, Valid: true},
		HashedPassword: hashedPassword,
	}

	err = s.store.UpdateUserHashedPassword(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

type forgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (s *Server) forgotPassword(ctx *gin.Context) {
	var req forgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	endpoint := fmt.Sprintf("https://%s/users/forgot-password", s.config.Domain)
	emailTemplate := filepath.Join(s.config.TemplatesDir, "password_reset.html")
	if err := SendPasswordResetMail(
		s.emailSender,
		user.ID.Bytes,
		user.Email,
		s.tokenMaker,
		s.config.AccessTokenDuration,
		endpoint, emailTemplate); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password reset email sent successfully"})
}

func SendPasswordResetMail(
	sender mail.EmailSender,
	userId uuid.UUID,
	email string,
	tokenMaker token.Maker,
	duration time.Duration,
	endpoint string,
	emailTemplate string,
) error {
	resetToken, _, err := tokenMaker.CreateToken(userId, duration, token.PurposeResetPassword)
	if err != nil {
		return err
	}

	passwordResetLink := fmt.Sprintf("%s?token=%s", endpoint, resetToken)
	tmp, err := template.ParseFiles(emailTemplate)
	if err != nil {
		return fmt.Errorf("parsing forgot password template failed: %w", err)
	}

	content := bytes.NewBufferString("")
	err = tmp.Execute(content, map[string]interface{}{
		"ResetURL": passwordResetLink,
	})
	if err != nil {
		return fmt.Errorf("executing forgot password template failed: %w", err)
	}

	subject := "Password Reset Request"
	return sender.SendMail(subject, content.String(), []string{email}, nil, nil, nil)
}
