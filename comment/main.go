package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rahulbharuka/github-proxy/comment/logic"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// set release mode logging.
	gin.SetMode(gin.ReleaseMode)

	// create default Gin router
	router := gin.New()

	// init log middleware
	router.Use(gin.Logger())

	// init recovery middleware
	router.Use(gin.Recovery())

	// get logic handler
	h := logic.GetHandler()

	// API handlers.
	router.POST("/orgs/:org/comments", h.PostComment)
	router.GET("/orgs/:org/comments", h.ListAllComments)
	router.DELETE("/orgs/:org/comments", h.DeleteAllComments)

	// run app on the specified port
	router.Run(":" + port)
}
