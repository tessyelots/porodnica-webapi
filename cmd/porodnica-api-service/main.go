package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tessyelots/porodnica-webapi/api"
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
	// request routings
	handleFunctions := &porodnica_ambulance_home.ApiHandleFunctions{
		PorodnicaWaitingListAPI: porodnica_ambulance_home.NewPorodnicaWaitingListApi(),
	}
	porodnica_ambulance_home.NewRouterWithGinEngine(engine, *handleFunctions)
	engine.GET("/openapi", api.HandleOpenApi)
	engine.Run(":" + port)
}
