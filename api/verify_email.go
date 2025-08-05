package api

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/sanjayj369/retrospect-backend/db/sqlc"
	mail "github.com/sanjayj369/retrospect-backend/mail"
	"github.com/sanjayj369/retrospect-backend/token"
)

type ResendVerificationEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (s *Server) ResendVerificationEmail(ctx *gin.Context) {
	var req ResendVerificationEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	if user.IsVerified {
		ctx.JSON(http.StatusOK, gin.H{"message": "Email already verified"})
		return
	}

	endpoint := fmt.Sprintf("https://%s/users/verify-email", s.config.Domain)
	emailTemplate := filepath.Join(s.config.TemplatesDir, "email_verification.html")
	if err := SendVerificationMail(
		s.emailSender,
		user.ID.Bytes,
		user.Email,
		s.tokenMaker,
		s.config.AccessTokenDuration,
		endpoint, emailTemplate); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Verification email sent successfully"})
}

func SendVerificationMail(
	sender mail.EmailSender,
	userId uuid.UUID,
	to string,
	tokenMaker token.Maker,
	duration time.Duration,
	endpoint string,
	templateFile string) error {
	tkn, _, err := tokenMaker.CreateToken(userId, duration, token.PurposeVerifyEmail)
	if err != nil {
		return fmt.Errorf("unable to create verification token: %w", err)
	}

	verficationLink := fmt.Sprintf("%s?token=%s", endpoint, tkn)
	tmp, err := template.ParseFiles(templateFile)
	if err != nil {
		return fmt.Errorf("parsing email template failed: %w", err)
	}

	content := bytes.NewBufferString("")
	err = tmp.Execute(content, map[string]string{
		"VerificationURL": verficationLink,
	})
	if err != nil {
		return fmt.Errorf("executing email template failed: %w", err)
	}

	subject := "Verify your email address"
	err = sender.SendMail(subject, content.String(), []string{to}, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("sending email failed: %w", err)
	}

	return nil
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

func (s *Server) VerifyEmail(ctx *gin.Context) {
	var req VerifyEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := s.tokenMaker.VerifyToken(req.Token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if payload.Purpose != token.PurposeVerifyEmail {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid token purpose")))
		return
	}

	arg := db.UpdateUserIsVerifiedParams{
		ID:         pgtype.UUID{Bytes: payload.UserId, Valid: true},
		IsVerified: true,
	}

	_, err = s.store.UpdateUserIsVerified(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}
