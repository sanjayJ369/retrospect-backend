package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/sanjayj369/retrospect-backend/db/sqlc"
)

type userChallengesRequest struct {
	ID       string `uri:"id" binding:"required,uuid"`
	PageIdx  int32  `form:"page_idx"`
	PageSize int32  `form:"page_size"`
}

func (server *Server) getUserChallenges(c *gin.Context) {
	var req userChallengesRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
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
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var userID [16]byte = id

	arg := db.ListChallengesByUserParams{
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
		Limit:  req.PageSize,
		Offset: req.PageIdx * req.PageSize,
	}

	challenges, err := server.store.ListChallengesByUser(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, challenges)
}
