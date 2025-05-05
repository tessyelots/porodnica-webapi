package porodnica_ambulance_home

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tessyelots/porodnica-webapi/internal/db_service"
)

type porodnicaUpdater = func(
	ctx *gin.Context,
	porodnica *Porodnica,
) (updatedPorodnica *Porodnica, responseContent interface{}, status int)

func updatePorodnicaFunc(ctx *gin.Context, updater porodnicaUpdater) {
	value, exists := ctx.Get("db_service")
	if !exists {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service not found",
				"error":   "db_service not found",
			})
		return
	}

	db, ok := value.(db_service.DbService[Porodnica])
	if !ok {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service context is not of type db_service.DbService",
				"error":   "cannot cast db_service context to db_service.DbService",
			})
		return
	}

	porodnicaId := ctx.Param("porodnicaId")

	porodnica, err := db.FindDocument(ctx, porodnicaId)

	switch err {
	case nil:
		// continue
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Porodnica not found",
				"error":   err.Error(),
			},
		)
		return
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to load porodnica from database",
				"error":   err.Error(),
			})
		return
	}

	if !ok {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "Failed to cast porodnica from database",
				"error":   "Failed to cast porodnica from database",
			})
		return
	}

	updatedPorodnica, responseObject, status := updater(ctx, porodnica)

	if updatedPorodnica != nil {
		err = db.UpdateDocument(ctx, porodnicaId, updatedPorodnica)
	} else {
		err = nil // redundant but for clarity
	}

	switch err {
	case nil:
		if responseObject != nil {
			ctx.JSON(status, responseObject)
		} else {
			ctx.AbortWithStatus(status)
		}
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Porodnica was deleted while processing the request",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to update porodnica in database",
				"error":   err.Error(),
			})
	}

}
