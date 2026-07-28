package main

import (
	"github.com/gin-gonic/gin"
	"go-api/controller"
)

func main(){
	server := gin.Default()

	//Camada Controller
	ProductController := controller.NewProductController()

	server.GET("/ping", func(ctx * gin.Context){
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	server.GET("/products", ProductController.GetProducts)
	server.Run(":8080")
}