package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"go-api/model"
	"go-api/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

type ProductUseCase struct {
	repository repository.ProductRepository
	redis      *redis.Client
	ctx        context.Context
}

func NewProductUseCase(repo repository.ProductRepository, redisClient *redis.Client) ProductUseCase {
	return ProductUseCase{
		repository: repo,
		redis:      redisClient,
		ctx:        context.Background(),
	}
}

func (pu *ProductUseCase) GetProducts() ([]model.Product, error) {
	cacheKey := "cache:products:all"

	// Tenta obter os produtos do cache Redis
	cached, err := pu.redis.Get(pu.ctx, cacheKey).Result()
	if err == nil {
		var products []model.Product

		if json.Unmarshal([]byte(cached), &products) == nil {
			return products, nil
		}
	}

	// Se não houver cache, busca os produtos no banco de dados
	products, err := pu.repository.GetProducts()
	if err != nil {
		return nil, err
	}

	// Armazena os produtos em cache no Redis por 60 segundos
	if b, mErr := json.Marshal(products); mErr == nil {
		_ = pu.redis.Set(pu.ctx, cacheKey, b, 60*time.Second).Err()
	}

	return products, nil
}

func (pu *ProductUseCase) CreateProduct(product model.Product) (model.Product, error) {
	productId, err := pu.repository.CreateProduct(product)

	if err != nil {
		return model.Product{}, err
	}

	product.Idq = productId
	pu.invalidateProductCache(productId)

	return product, nil
}

func (pu *ProductUseCase) GetProductById(id_product int) (*model.Product, error) {
	cacheKey := fmt.Sprintf("cache:products:%d", id_product)

	// Tenta obter o produto do cache Redis
	cached, err := pu.redis.Get(pu.ctx, cacheKey).Result()
	if err == nil {
		var product model.Product

		if json.Unmarshal([]byte(cached), &product) == nil {
			return &product, nil
		}
	}

	// Se não houver cache, busca o produto no banco de dados
	product, err := pu.repository.GetProductById(id_product)

	if err != nil {
		return nil, err
	}

	if b, mErr := json.Marshal(product); mErr == nil {
		_ = pu.redis.Set(pu.ctx, cacheKey, b, 60*time.Second).Err()
	}

	return product, nil
}

func (pu *ProductUseCase) UpdateProduct(id_product int, product model.ProductRequest) (*model.Product, error) {
	updatedProduct, err := pu.repository.UpdateProduct(id_product, product)
	if err != nil {
		return nil, err
	}

	pu.invalidateProductCache(id_product)
	return updatedProduct, nil
}

func (pu *ProductUseCase) DeleteProduct(id_product int) (bool, error) {
	isDeleted, err := pu.repository.DeleteProduct(id_product)
	if err != nil {
		return false, err
	}

	if isDeleted {
		pu.invalidateProductCache(id_product)
	}

	return isDeleted, nil
}

func (pu *ProductUseCase) invalidateProductCache(id int) {
	_ = pu.redis.Del(
		pu.ctx,
		"cache:products:all",
		fmt.Sprintf("cache:products:%d", id),
		fmt.Sprintf("cache:products:id:%d", id),
	).Err()
}
