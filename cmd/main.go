package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go-api/controller"
	"go-api/db"
	"go-api/repository"
	"go-api/usecase"
)

func main() {
	server := gin.Default()

	dbConnect, err := db.ConectDB()

	if err != nil {
		panic(err)
	}

	//Camada de Repository
	ProductRepository := repository.NewProductRepository(dbConnect)

	//Camadade UseCase
	ProductUseCase := usecase.NewProductUseCase(ProductRepository)

	//Camada Controller
	ProductController := controller.NewProductController(ProductUseCase)

	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//Rotas
	server.GET("/products", ProductController.GetProducts)
	server.GET("/products/:id", ProductController.GetProductById)
	server.POST("/products", ProductController.CreateProduct)
	server.PUT("/products/:id", ProductController.UpdateProduct)
	server.DELETE("/products/:id", ProductController.DeleteProduct)
	server.Run(":8080")
}
