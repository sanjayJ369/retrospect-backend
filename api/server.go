package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/sanjayj369/retrospect-backend/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// user routes
	router.POST("/users", server.createUser)
	router.GET("/users/:id", server.getUser)

	// challenge routes
	router.POST("/challenges", server.createChallenge)
	router.GET("/challenges/:id", server.getChallenge) // TODO: update to include challenge entries
	router.PUT("/challenges/:id", server.updateChallenge)
	router.DELETE("/challenges/:id", server.deleteChallenge)
	router.GET("/challenges", server.listChallenges)

	// challenge entries
	router.PUT("/challenge-entries/:id", server.updateChallengeEntries)

	server.router = router
	return server
}

// Start runs server on provided address
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

// errorResponse converts error to a json object
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
