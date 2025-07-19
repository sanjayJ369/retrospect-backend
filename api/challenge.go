package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/sanjayj369/retrospect-backend/db/sqlc"
)

type createChallengeRequest struct {
	Title       string `json:"title" binding:"required"`
	UserID      string `json:"user_id" binding:"required"`
	Description string `json:"description"`
	EndDate     string `json:"end_date"`
}

func (server *Server) createChallenge(c *gin.Context) {
	var req createChallengeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, errorResponse(err))
		return
	}

	parseUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var userid [16]byte = parseUUID
	var endDatePg pgtype.Date
	endDate, err := time.Parse("2006-01-02", req.EndDate)

	if err != nil {
		endDatePg.Valid = false
	} else {
		endDatePg = pgtype.Date{Time: endDate, Valid: true}
	}

	arg := db.CreateChallengeParams{
		Title:       req.Title,
		UserID:      pgtype.UUID{Bytes: userid, Valid: true},
		Description: pgtype.Text{String: req.Description, Valid: true},
		EndDate:     endDatePg,
	}

	challenge, err := server.store.CreateChallenge(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, challenge)
}

type getChallengeRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func (server *Server) getChallenge(c *gin.Context) {
	var req getChallengeRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var challengeID [16]byte = id
	challenge, err := server.store.GetChallenge(c, pgtype.UUID{Bytes: challengeID, Valid: true})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, challenge)
}

type updateChallengeRequest struct {
	ID          string `uri:"id" binding:"required,uuid"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	EndDate     string `json:"end_date"`
}

func (server *Server) updateChallenge(c *gin.Context) {
	var req updateChallengeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var challengeID [16]byte = id

	var endDatePg pgtype.Date
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		endDatePg.Valid = false
	} else {
		endDatePg = pgtype.Date{Time: endDate, Valid: true}
	}

	arg := db.UpdateChallengeDetailsParams{
		ID:          pgtype.UUID{Bytes: challengeID, Valid: true},
		Title:       req.Title,
		Description: pgtype.Text{String: req.Description, Valid: true},
		EndDate:     endDatePg,
	}

	challenge, err := server.store.UpdateChallengeDetails(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, challenge)
}

type deleteChallengeRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func (server *Server) deleteChallenge(c *gin.Context) {
	var req deleteChallengeRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var challengeID [16]byte = id

	challenge, err := server.store.DeleteChallenge(c, pgtype.UUID{Bytes: challengeID, Valid: true})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, challenge)
}
