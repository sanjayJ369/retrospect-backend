package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/sanjayj369/retrospect-backend/db/sqlc"
	"github.com/sanjayj369/retrospect-backend/mail"
	"github.com/sanjayj369/retrospect-backend/ratelimiter"
	"github.com/sanjayj369/retrospect-backend/token"
	"github.com/sanjayj369/retrospect-backend/util"
)

type Server struct {
	config      util.Config
	store       db.Store
	tokenMaker  token.Maker
	router      *gin.Engine
	emailSender mail.EmailSender
	rateLimiter ratelimiter.RateLimiter
}

func NewServer(config util.Config, store db.Store, emailSender mail.EmailSender, rateLimiter ratelimiter.RateLimiter) (*Server, error) {
	maker, err := token.NewPasetoMaker(config.SymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:      config,
		store:       store,
		tokenMaker:  maker,
		emailSender: emailSender,
		rateLimiter: rateLimiter,
	}

	setupRoutes(server)

	return server, nil
}

// Start runs server on provided address
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

// errorResponse converts error to a json object
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func setupRoutes(server *Server) {
	router := gin.Default()

	// user routes
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.LoginUser)
	router.POST("/tokens/renew_access", server.RenewAccessToken)

	// rate limiting middleware
	rateLimited := router.Group("/").Use(ratelimiterMiddleware(
		server.rateLimiter,
		server.config.RatelimitDuration))
	rateLimited.POST("/users/verify-email", server.VerifyEmail)
	rateLimited.POST("/users/resend-verification", server.ResendVerificationEmail)

	// authorized routers
	authRouter := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRouter.GET("/users/:id", server.getUser)

	// challenge routes
	authRouter.POST("/challenges", server.createChallenge)
	authRouter.GET("/challenges/:id", server.getChallenge) // TODO: update to include challenge entries
	authRouter.PATCH("/challenges/:id", server.updateChallenge)
	authRouter.DELETE("/challenges/:id", server.deleteChallenge)
	authRouter.GET("/users/:id/challenges", server.listChallenges)

	// challenge entries
	authRouter.PUT("/challenge-entries/:id", server.updateChallengeEntries)

	// tasks
	authRouter.POST("/tasks", server.createTask)
	authRouter.GET("/tasks/:id", server.getTask)
	authRouter.PATCH("/tasks/:id", server.updateTask)
	authRouter.DELETE("/tasks/:id", server.deleteTask)
	authRouter.GET("/tasks", server.listTasks)

	server.router = router
}
