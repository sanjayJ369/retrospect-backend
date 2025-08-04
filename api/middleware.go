package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sanjayj369/retrospect-backend/ratelimiter"
	"github.com/sanjayj369/retrospect-backend/token"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func ratelimiterMiddleware(limiter ratelimiter.RateLimiter, duration time.Duration) gin.HandlerFunc {
	type emailRequest struct {
		Email string `json:"email" binding:"email"`
	}

	return func(ctx *gin.Context) {
		key := ctx.ClientIP()

		allowed, err := limiter.Allow(key, duration)
		if err != nil || !allowed {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, errorResponse(fmt.Errorf("rate limit exceeded")))
			return
		}

		var bodyBytes []byte
		if ctx.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(ctx.Request.Body)
		}

		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var req emailRequest
		if err := json.Unmarshal(bodyBytes, &req); err == nil {
			if req.Email != "" {
				allowed, err := limiter.Allow(req.Email, duration)
				if err != nil {
					ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("internal server error")))
					return
				}
				if !allowed {
					ctx.AbortWithStatusJSON(http.StatusTooManyRequests, errorResponse(fmt.Errorf("rate limit exceeded for this email")))
					return
				}
			}
		}

		ctx.Next()
	}
}

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := fmt.Errorf("authorization header is not present")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) != 2 {
			err := fmt.Errorf("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type: %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
