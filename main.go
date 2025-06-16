package main

import (
	"fmt"
	"log"
	"os"
	"pasteGo/backend/api/rest/middlewares"
	"pasteGo/backend/api/rest/v1/handlers"
	"pasteGo/backend/api/rest/v1/types"
	"pasteGo/backend/db"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	importENV()
	dbInstance, err := db.GetDBInstance()
	defer db.CloseDB()
	if err != nil {
		log.Fatalf("Ошибка при инициализации базы данных: %s", err)
	}

	err = dbInstance.Init()
	if err != nil {
		log.Fatalf("Ошибка при инициализации базы данных: %s", err)
	}

	router := gin.Default()

	router.Static("/_app/immutable/", "./build/_app/immutable/")
	router.NoRoute(func(c *gin.Context) {
		c.File("./build/index.html")
	})

	//! dev headers
	router.Use(middlewares.CORSMiddleware())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	rest := router.Group("/rest")
	{
		rest.POST("/auth", handlers.Login)
		rest.POST("/registration", handlers.Register)
		rest.DELETE("/logout", handlers.Logout)
		rest.POST("/update_tokens", middlewares.JwtRefreshMiddleware(), handlers.Refresh)

		rest.POST("/paste/:id", handlers.GetPaste)

		v1 := rest.Group("/v1", middlewares.JwtMiddleware())
		{
			v1.GET("/testtoken", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "pong",
				})
			})

			v1.PUT("/user", handlers.UpdateUser)
			v1.DELETE("/user", handlers.DeleteUser)

			v1.GET("/paste", handlers.GetPasteList)
			v1.POST("/paste", handlers.CreatePaste)
			v1.PUT("/paste/:id", handlers.UpdatePaste)
			v1.DELETE("/paste/:id", handlers.DeletePaste)
		}
	}

	router.Run("0.0.0.0:10015")
}

func importENV() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Ошибка при парсинге .env файла: %s", err)
	}

	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		log.Fatalf("SECRET_KEY не установлен")
	}

	types.SecretKey = []byte(secret)
	fmt.Println(secret)
}
