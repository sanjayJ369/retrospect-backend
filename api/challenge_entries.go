package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/sanjayj369/retrospect-backend/db/sqlc"
)

type updateChallengeEntriesUriRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type updateChallengeEntriesBodyRequest struct {
	Complete bool `json:"complete"`
}

func (server *Server) updateChallengeEntries(ctx *gin.Context) {
	var uriReq updateChallengeEntriesUriRequest
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var bodyReq updateChallengeEntriesBodyRequest
	if err := ctx.ShouldBindJSON(&bodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	parsedEntryUUID, _ := uuid.Parse(uriReq.ID)
	entryID := pgtype.UUID{Bytes: parsedEntryUUID, Valid: true}

	entry, err := server.store.GetChallengeEntry(ctx, entryID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	challenge, err := server.store.GetChallenge(ctx, entry.ChallengeID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	if err := authorizeUser(ctx, challenge.UserID.Bytes); err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	arg := db.UpdateChallengeEntryParams{
		ID:        entryID,
		Completed: pgtype.Bool{Bool: bodyReq.Complete, Valid: true},
	}
	res, err := server.store.UpdateChallengeEntry(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}
