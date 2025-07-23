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

func (server *Server) createChallenge(ctx *gin.Context) {
	var req createChallengeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	parseUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
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

	err = authorizeUser(ctx, userid)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	arg := db.CreateChallengeParams{
		Title:       req.Title,
		UserID:      pgtype.UUID{Bytes: userid, Valid: true},
		Description: pgtype.Text{String: req.Description, Valid: true},
		EndDate:     endDatePg,
	}

	challenge, err := server.store.CreateChallenge(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, challenge)
}

type getChallengeRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func (server *Server) getChallenge(ctx *gin.Context) {
	// TODO: add challenge entires along with challenge details
	var req getChallengeRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var challengeID [16]byte = id
	challenge, err := server.store.GetChallenge(ctx, pgtype.UUID{Bytes: challengeID, Valid: true})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = authorizeUser(ctx, challenge.UserID.Bytes)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, challenge)
}

type updateChallengeRequest struct {
	ID          string `uri:"id" binding:"required,uuid"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	EndDate     string `json:"end_date"`
}

func (server *Server) updateChallenge(ctx *gin.Context) {
	var req updateChallengeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
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

	challenge, err := server.store.GetChallenge(ctx, arg.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if err := authorizeUser(ctx, challenge.UserID.Bytes); err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	newChallenge, err := server.store.UpdateChallengeDetails(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newChallenge)
}

type deleteChallengeRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func (server *Server) deleteChallenge(ctx *gin.Context) {
	var req deleteChallengeRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var challengeID [16]byte = id
	challengeIDPGType := pgtype.UUID{Bytes: challengeID, Valid: true}

	challenge, err := server.store.GetChallenge(ctx, challengeIDPGType)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	if err := authorizeUser(ctx, challenge.UserID.Bytes); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	deletedChallenge, err := server.store.DeleteChallenge(ctx, challengeIDPGType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, deletedChallenge)
}

type userChallengesRequest struct {
	ID       string `uri:"id" binding:"required,uuid"`
	PageIdx  int32  `form:"page_idx"`
	PageSize int32  `form:"page_size"`
}

func (server *Server) listChallenges(ctx *gin.Context) {
	var req userChallengesRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageIdx < 0 {
		req.PageIdx = 0
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var userID [16]byte = id

	if err := authorizeUser(ctx, userID); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.ListChallengesByUserParams{
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
		Limit:  req.PageSize,
		Offset: req.PageIdx * req.PageSize,
	}

	challenges, err := server.store.ListChallengesByUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, challenges)
}
