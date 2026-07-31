package repository

import (
	"database/sql"
	"fmt"
	"go-api/model"
)

type ProductRepository struct {
	connection *sql.DB
}

func NewProductRepository(connection *sql.DB) ProductRepository {
	return ProductRepository{
		connection: connection,
	}
}

func (pr *ProductRepository) GetProducts() ([]model.Product, error) {
	var query = "SELECT id, product_name, price FROM products WHERE excluded_bit = 0::BIT"
	rows, err := pr.connection.Query(query)

	if err != nil {
		fmt.Print(err)
		return []model.Product{}, err
	}

	var productList []model.Product
	var productObj model.Product

	for rows.Next() {
		err = rows.Scan(
			&productObj.Idq,
			&productObj.Name,
			&productObj.Price)

		if err != nil {
			fmt.Print(err)
			return []model.Product{}, err
		}

		productList = append(productList, productObj)
	}
	rows.Close()

	return productList, nil
}

func (pr *ProductRepository) CreateProduct(product model.Product) (int, error) {
	var id int
	query, err := pr.connection.Prepare("INSERT INTO products " +
		"(product_name, price) " +
		"VALUES ($1, $2) RETURNING id")

	if err != nil {
		fmt.Print(err)
		return 0, err
	}
	defer query.Close()

	err = query.QueryRow(product.Name, product.Price).Scan(&id)
	if err != nil {
		fmt.Print(err)
		return 0, err
	}

	return id, nil
}

func (pr *ProductRepository) GetProductById(id_product int) (*model.Product, error) {
	query, err := pr.connection.Prepare("SELECT id, product_name, price FROM products WHERE id = $1 AND excluded_bit = 0::BIT")

	if err != nil {
		print(err)
		return nil, err
	}

	var product model.Product
	err = query.QueryRow(id_product).Scan(
		&product.Idq,
		&product.Name,
		&product.Price,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	query.Close()
	return &product, nil
}

func (pr *ProductRepository) UpdateProduct(id_product int, product model.ProductRequest) (*model.Product, error) {
	query, err := pr.connection.Prepare(
		"UPDATE products SET product_name = $1, price = $2 WHERE id = $3 AND excluded_bit = 0::BIT RETURNING id, product_name, price",
	)

	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	defer query.Close()

	var updatedProduct model.Product
	err = query.QueryRow(product.Name, product.Price, id_product).Scan(
		&updatedProduct.Idq,
		&updatedProduct.Name,
		&updatedProduct.Price,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		fmt.Print(err)
		return nil, err
	}

	return &updatedProduct, nil
}

func (pr *ProductRepository) DeleteProduct(id_product int) (bool, error) {
	query, err := pr.connection.Prepare(
		"UPDATE products SET excluded_bit = 1::BIT WHERE id = $1 AND excluded_bit = 0::BIT",
	)
	if err != nil {
		fmt.Print(err)
		return false, err
	}
	defer query.Close()

	result, err := query.Exec(id_product)
	if err != nil {
		fmt.Print(err)
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Print(err)
		return false, err
	}

	return rowsAffected > 0, nil
}
