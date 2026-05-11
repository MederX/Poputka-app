package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"kyrgyzstan-rideshare/config"
	"kyrgyzstan-rideshare/firebase"
	"kyrgyzstan-rideshare/handlers"
	"kyrgyzstan-rideshare/middleware"
)

func main() {
	godotenv.Load()
	config.Load()
	firebase.Init()

	r := gin.Default()
	r.Use(corsMiddleware())

	r.POST("/auth/telegram", handlers.AuthTelegram)

	auth := r.Group("/", middleware.AuthRequired(config.C.JWTSecret))
	{
		auth.GET("/users/me", handlers.GetMe)
		auth.PUT("/users/me", handlers.UpdateMe)

		auth.GET("/rides", handlers.ListRides)
		auth.POST("/rides", handlers.CreateRide)
		auth.GET("/rides/:id", handlers.GetRide)
		auth.PUT("/rides/:id", handlers.UpdateRide)
		auth.DELETE("/rides/:id", handlers.DeleteRide)
		auth.PATCH("/rides/:id/complete", handlers.CompleteRide)

		auth.POST("/rides/:id/requests", handlers.CreateRideRequest)
		auth.GET("/rides/:id/requests", handlers.ListRideRequests)
		auth.PATCH("/rides/:id/requests/:reqId", handlers.UpdateRideRequest)

		auth.GET("/passenger-posts", handlers.ListPassengerPosts)
		auth.POST("/passenger-posts", handlers.CreatePassengerPost)
		auth.DELETE("/passenger-posts/:id", handlers.DeletePassengerPost)

		auth.GET("/cargo", handlers.ListCargo)
		auth.POST("/cargo", handlers.CreateCargo)
		auth.DELETE("/cargo/:id", handlers.DeleteCargo)
	}

	port := config.C.Port
	log.Printf("Starting server on :%s", port)
	log.Fatal(r.Run(":" + port))
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type,Accept,Origin,User-Agent,DNT,Cache-Control,X-Mx-ReqToken,Keep-Alive,X-Requested-With,If-Modified-Since")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
