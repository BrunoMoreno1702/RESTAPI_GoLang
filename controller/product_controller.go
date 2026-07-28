package controller

import (
	"go-api/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type productController struct {
	//Usecase
}

func NewProductController() productController {
	return productController{}
} 

func (p *productController) GetProducts (ctx *gin.Context){
	//Mock de dado para teste
	products := []model.Product{
		{
			Idq: 1, 
			Name: "Batata Frita", 
			Price: 10.0,
		},
	}
	
	ctx.JSON(http.StatusOK, products)
}