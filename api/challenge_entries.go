package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/sanjayj369/retrospect-backend/db/sqlc"
)

type updateChallengeEntriesUriRequest struct {
	ChallengeID string `uri:"id" binding:"required,uuid"`
}

type updateChallengeEntriesBodyRequest struct {
	Complete bool `json:"complete"`
}

func (server *Server) updateChallengeEntries(ctx *gin.Context) {
	var uriReq updateChallengeEntriesUriRequest
	err := ctx.ShouldBindUri(&uriReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	parseUUID, err := uuid.Parse(uriReq.ChallengeID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var challengeID [16]byte = parseUUID

	var bodyReq updateChallengeEntriesBodyRequest
	err = ctx.ShouldBindJSON(&bodyReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateChallengeEntryParams{
		ID:        pgtype.UUID{Bytes: challengeID, Valid: true},
		Completed: pgtype.Bool{Bool: bodyReq.Complete, Valid: true},
	}

	res, err := server.store.UpdateChallengeEntry(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}
