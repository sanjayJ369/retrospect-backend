package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/sanjayj369/retrospect-backend/db/sqlc"
	"github.com/sanjayj369/retrospect-backend/util"
)

type createTaskRequest struct {
	TaskDayID   string `json:"task_day_id" binding:"required,uuid"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Duration    int    `json:"duration" binding:"required"`
}

func (s *Server) createTask(ctx *gin.Context) {
	var req createTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	parseUUID, err := uuid.Parse(req.TaskDayID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	taskDayID := pgtype.UUID{Bytes: parseUUID, Valid: true}
	description := pgtype.Text{String: req.Description, Valid: req.Description != ""}
	duration := util.MinutesToPGInterval(req.Duration)

	arg := db.CreateTaskParams{
		TaskDayID:   taskDayID,
		Title:       req.Title,
		Description: description,
		Duration:    duration,
	}

	task, err := s.store.CreateTask(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, task)
}

type getTaskRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func (s *Server) getTask(ctx *gin.Context) {
	var req getTaskRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	parasedUUID, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var taskID [16]byte = parasedUUID

	task, err := s.store.GetTask(ctx, pgtype.UUID{Bytes: taskID, Valid: true})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, task)
}

type updateTaskUriRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type updateTaskBodyRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	Duration    int    `json:"duration" binding:"required"`
}

func (s *Server) updateTask(ctx *gin.Context) {
	var uriReq updateTaskUriRequest
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var bodyReq updateTaskBodyRequest
	if err := ctx.ShouldBindJSON(&bodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	parseUUID, err := uuid.Parse(uriReq.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	taskID := pgtype.UUID{Bytes: parseUUID, Valid: true}
	description := pgtype.Text{String: bodyReq.Description, Valid: bodyReq.Description != ""}
	duration := util.MinutesToPGInterval(bodyReq.Duration)
	completed := pgtype.Bool{Bool: bodyReq.Completed, Valid: true}

	arg := db.UpdateTaskParams{
		ID:          taskID,
		Title:       bodyReq.Title,
		Description: description,
		Completed:   completed,
		Duration:    duration,
	}

	res, err := s.store.UpdateTask(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}

type deleteTaskRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func (s *Server) deleteTask(ctx *gin.Context) {
	var req deleteTaskRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	parseUUID, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	taskID := pgtype.UUID{Bytes: parseUUID, Valid: true}

	res, err := s.store.DeleteTask(ctx, taskID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, res)
}

type listTasksRequest struct {
	TaskDayID string `form:"task_day_id" binding:"required,uuid"`
}

func (s *Server) listTasks(ctx *gin.Context) {
	var req listTasksRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	parseUUID, err := uuid.Parse(req.TaskDayID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	taskDayID := pgtype.UUID{Bytes: parseUUID, Valid: true}

	tasks, err := s.store.ListTasksByTaskDayId(ctx, taskDayID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}
