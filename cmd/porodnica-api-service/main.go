package main

import (
	"log"
	"os"
	"strings"

	"context"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tessyelots/porodnica-webapi/api"
	"github.com/tessyelots/porodnica-webapi/internal/db_service"
	"github.com/tessyelots/porodnica-webapi/internal/porodnica_ambulance_home"
)

func main() {
	log.Printf("Server started")
	port := os.Getenv("AMBULANCE_API_PORT")
	if port == "" {
		port = "8080"
	}
	environment := os.Getenv("AMBULANCE_API_ENVIRONMENT")
	if !strings.EqualFold(environment, "production") { // case insensitive comparison
		gin.SetMode(gin.DebugMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())
	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{""},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})
	engine.Use(corsMiddleware)

	// setup context update  middleware
	dbService := db_service.NewMongoService[porodnica_ambulance_home.Porodnica](db_service.MongoServiceConfig{})
	defer dbService.Disconnect(context.Background())
	engine.Use(func(ctx *gin.Context) {
		ctx.Set("db_service", dbService)
		ctx.Next()
	})
	// request routings
	handleFunctions := &porodnica_ambulance_home.ApiHandleFunctions{
		PorodnicaWaitingListAPI: porodnica_ambulance_home.NewPorodnicaWaitingListApi(),
		PorodniceAPI:            porodnica_ambulance_home.NewPorodniceApi(),
	}
	porodnica_ambulance_home.NewRouterWithGinEngine(engine, *handleFunctions)
	engine.GET("/openapi", api.HandleOpenApi)
	engine.Run(":" + port)
}
