package main

import (
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go-api/controller"
	"go-api/db"
	"go-api/middleware"
	"go-api/repository"
	"go-api/usecase"
)

func main() {
	server := gin.Default()

	dbConnect, err := db.ConectDB()
	if err != nil {
		panic(err)
	}

	redisClient, err := db.ConnectRedis()
	if err != nil {
		panic(err)
	}

	// Camada de Repository
	ProductRepository := repository.NewProductRepository(dbConnect)

	// Camada de UseCase
	ProductUseCase := usecase.NewProductUseCase(ProductRepository, redisClient)

	// Camada Controller
	ProductController := controller.NewProductController(ProductUseCase)

	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	products := server.Group("/products")
	products.Use(middleware.RateLimit(redisClient, 60, time.Minute))
	{
		products.GET("", ProductController.GetProducts)
		products.GET("/:id", ProductController.GetProductById)
		products.POST("", ProductController.CreateProduct)
		products.PUT("/:id", ProductController.UpdateProduct)
		products.DELETE("/:id", ProductController.DeleteProduct)
	}

	server.Run(":8080")
}
