package porodnica_ambulance_home

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tessyelots/porodnica-webapi/internal/db_service"
)

type implPorodniceAPI struct {
}

func NewPorodniceApi() PorodniceAPI {
	return &implPorodniceAPI{}
}

func (o implPorodniceAPI) CreatePorodnica(c *gin.Context) {
	// get db service from context
	value, exists := c.Get("db_service")
	if !exists {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db not found",
				"error":   "db not found",
			})
		return
	}

	db, ok := value.(db_service.DbService[Porodnica])
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db context is not of required type",
				"error":   "cannot cast db context to db_service.DbService",
			})
		return
	}

	porodnica := Porodnica{}
	err := c.BindJSON(&porodnica)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  "Bad Request",
				"message": "Invalid request body",
				"error":   err.Error(),
			})
		return
	}

	if porodnica.Id == "" {
		porodnica.Id = uuid.New().String()
	}

	err = db.CreateDocument(c, porodnica.Id, &porodnica)

	switch err {
	case nil:
		c.JSON(
			http.StatusCreated,
			porodnica,
		)
	case db_service.ErrConflict:
		c.JSON(
			http.StatusConflict,
			gin.H{
				"status":  "Conflict",
				"message": "Porodnica already exists",
				"error":   err.Error(),
			},
		)
	default:
		c.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to create porodnica in database",
				"error":   err.Error(),
			},
		)
	}
}

func (o implPorodniceAPI) DeletePorodnica(c *gin.Context) {
	// get db service from context
	value, exists := c.Get("db_service")
	if !exists {
		c.JSON(
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
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service context is not of type db_service.DbService",
				"error":   "cannot cast db_service context to db_service.DbService",
			})
		return
	}

	porodnicaId := c.Param("porodnicaId")
	err := db.DeleteDocument(c, porodnicaId)

	switch err {
	case nil:
		c.AbortWithStatus(http.StatusNoContent)
	case db_service.ErrNotFound:
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Porodnica not found",
				"error":   err.Error(),
			},
		)
	default:
		c.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to delete ambulance from database",
				"error":   err.Error(),
			})
	}
}
