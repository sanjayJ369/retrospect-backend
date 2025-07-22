package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/sanjayj369/retrospect-backend/db/sqlc"
	"github.com/sanjayj369/retrospect-backend/token"
	"github.com/sanjayj369/retrospect-backend/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	maker, err := token.NewPasetoMaker(config.SymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: maker,
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
	router.GET("/users/:id", server.getUser)
	router.POST("/users/login", server.LoginUser)

	// challenge routes
	router.POST("/challenges", server.createChallenge)
	router.GET("/challenges/:id", server.getChallenge) // TODO: update to include challenge entries
	router.PATCH("/challenges/:id", server.updateChallenge)
	router.DELETE("/challenges/:id", server.deleteChallenge)
	router.GET("/users/:id/challenges", server.listChallenges)

	// challenge entries
	router.PUT("/challenge-entries/:id", server.updateChallengeEntries)

	// tasks
	router.POST("/tasks", server.createTask)
	router.GET("/tasks/:id", server.getTask)
	router.PATCH("/tasks/:id", server.updateTask)
	router.DELETE("/tasks/:id", server.deleteTask)
	router.GET("/tasks", server.listTasks)

	server.router = router
}
