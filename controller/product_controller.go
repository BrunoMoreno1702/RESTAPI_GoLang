package controller

import (
	"go-api/model"
	"go-api/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type productController struct {
	productUseCase usecase.ProductUseCase
}

func NewProductController(usecase usecase.ProductUseCase) productController {
	return productController{
		productUseCase: usecase,
	}
}

func (p *productController) GetProducts(ctx *gin.Context) {
	products, err := p.productUseCase.GetProducts()

	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			model.ErrorResponse("INTERNAL_ERROR", "Erro interno ao buscar produtos"),
		)
		return
	}

	ctx.JSON(http.StatusOK, model.SuccessResponse(products, "Produtos listados com sucesso"))
}

func (p *productController) CreateProduct(ctx *gin.Context) {
	var productRequest model.ProductRequest
	err := ctx.BindJSON(&productRequest)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			model.ErrorResponse("INVALID_BODY", "Payload inválido"),
		)
		return
	}

	product := model.Product{
		Name:  productRequest.Name,
		Price: productRequest.Price,
	}

	insertedProduct, err := p.productUseCase.CreateProduct(product)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			model.ErrorResponse("INTERNAL_ERROR", "Erro interno ao criar produto"),
		)
		return
	}

	ctx.JSON(http.StatusCreated, model.SuccessResponse(insertedProduct, "Produto criado com sucesso"))
}

func (p *productController) GetProductById(ctx *gin.Context) {

	id := ctx.Param("id")

	//Validações de entrada
	if id == "" {
		ctx.JSON(
			http.StatusBadRequest,
			model.ErrorResponse("INVALID_ID", "Id do produto não pode ser nulo ou vazio"),
		)
		return
	}

	productId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			model.ErrorResponse("INVALID_ID", "Id do produto precisa ser um numero inteiro"),
		)
		return
	}

	products, err := p.productUseCase.GetProductById(productId)

	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			model.ErrorResponse("INTERNAL_ERROR", "Erro interno ao buscar produto"),
		)
		return
	}

	if products == nil {
		ctx.JSON(
			http.StatusNotFound,
			model.ErrorResponse("PRODUCT_NOT_FOUND", "Produto não encontrado"),
		)
		return

	}
	ctx.JSON(http.StatusOK, model.SuccessResponse(products, "Produto encontrado"))
}

func (p *productController) UpdateProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(
			http.StatusBadRequest,
			model.ErrorResponse("INVALID_ID", "Id do produto não pode ser nulo ou vazio"),
		)
		return
	}

	productId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			model.ErrorResponse("INVALID_ID", "Id do produto precisa ser um numero inteiro"),
		)
		return
	}

	var product model.ProductRequest
	err = ctx.BindJSON(&product)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			model.ErrorResponse("INVALID_BODY", "Payload inválido"),
		)
		return
	}

	updatedProduct, err := p.productUseCase.UpdateProduct(productId, product)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			model.ErrorResponse("INTERNAL_ERROR", "Erro interno ao atualizar produto"),
		)
		return
	}

	if updatedProduct == nil {
		ctx.JSON(
			http.StatusNotFound,
			model.ErrorResponse("PRODUCT_NOT_FOUND", "Produto não encontrado"),
		)
		return
	}

	ctx.JSON(http.StatusOK, model.SuccessResponse(updatedProduct, "Produto atualizado com sucesso"))
}

func (p *productController) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(
			http.StatusBadRequest,
			model.ErrorResponse("INVALID_ID", "Id do produto não pode ser nulo ou vazio"),
		)
		return
	}

	productId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			model.ErrorResponse("INVALID_ID", "Id do produto precisa ser um numero inteiro"),
		)
		return
	}

	isDeleted, err := p.productUseCase.DeleteProduct(productId)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			model.ErrorResponse("INTERNAL_ERROR", "Erro interno ao excluir produto"),
		)
		return
	}

	if !isDeleted {
		ctx.JSON(
			http.StatusNotFound,
			model.ErrorResponse("PRODUCT_NOT_FOUND", "Produto não encontrado"),
		)
		return
	}

	ctx.Status(http.StatusNoContent)
}
