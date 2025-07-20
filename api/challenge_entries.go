package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/sanjayj369/retrospect-backend/db/sqlc"
)

type updateChallengeEntriesRequest struct {
	ChallengeID string `uri:"id" binding:"required,uuid"`
	Complete    bool   `json:"complete" binding:"required"`
}

func (server *Server) updateChallengeEntries(ctx *gin.Context) {
	var req updateChallengeEntriesRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	parseUUID, err := uuid.Parse(req.ChallengeID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var challengeID [16]byte = parseUUID

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateChallengeEntryParams{
		ID:        pgtype.UUID{Bytes: challengeID, Valid: true},
		Completed: pgtype.Bool{Bool: req.Complete, Valid: true},
	}

	res, err := server.store.UpdateChallengeEntry(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}
