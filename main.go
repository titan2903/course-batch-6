package main

import (
	"course-batch-6/internal/app/database"
	exerciseHandler "course-batch-6/internal/app/exercise/handler"

	userHandler "course-batch-6/internal/app/user/handler"
	"course-batch-6/internal/pkg/middleware"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{
			"message": "Successfully Deployed application",
		})
	})

	db := database.NewConnDatabasePostgres()
	exerciseHandler := exerciseHandler.NewExerciseHandler(db)
	userHandler := userHandler.NewUserHandler(db)
	r.POST("/exercises", middleware.WithAuh(), exerciseHandler.CreateExercise) //! Create exercise
	groupExerciseId := r.Group("/exercises/:exerciseId")
	{
		groupExerciseId.GET("", middleware.WithAuh(), exerciseHandler.GetExerciseByID)
		groupExerciseId.GET("score", middleware.WithAuh(), exerciseHandler.GetScore)
		groupExerciseId.POST("questions", middleware.WithAuh(), exerciseHandler.CreateQuestion)                  //! Create questions of the exercise
		groupExerciseId.POST("questions/:questionId/answer", middleware.WithAuh(), exerciseHandler.CreateAnswer) //! Create answer the question of the exercises
	}

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	r.Run(":1234")
}
