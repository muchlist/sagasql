package dao

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/muchlist/sagasql/db"
	"github.com/muchlist/sagasql/dto"
	"github.com/muchlist/sagasql/utils/rest_err"
)

func NewProductDao() ProductDaoAssumer {
	return &productDao{}
}

type ProductDaoAssumer interface {
	Insert(product dto.Product) (*int64, rest_err.APIError)
	Get(productID int64) (*dto.Product, rest_err.APIError)
	Find() ([]dto.Product, rest_err.APIError)
	Edit(productInput dto.Product) (*dto.Product, rest_err.APIError)
	Delete(productID int64) rest_err.APIError
}

type productDao struct {
}

func (u *productDao) Insert(product dto.Product) (*int64, rest_err.APIError) {
	sqlStatement := `
	INSERT INTO products (name, price, created_by, created_at) 
	VALUES ($1, $2, $3, $4) RETURNING product_id;
	`
	var productID int64
	err := db.DB.QueryRow(context.Background(), sqlStatement, product.Name, product.Price, product.CreatedBy, product.CreatedAt).Scan(&productID)
	if err != nil {
		return nil, rest_err.NewInternalServerError("gagal menambahkan product", err)
	}
	return &productID, nil
}

func (u *productDao) Edit(input dto.Product) (*dto.Product, rest_err.APIError) {
	sqlStatement := `
	UPDATE products 
	SET name = $2, price = $3
	WHERE product_id = $1 
	RETURNING product_id, name, price, created_by, created_at;
	`

	var product dto.Product
	err := db.DB.QueryRow(
		context.Background(),
		sqlStatement, input.ProductID, input.Name, input.Price,
	).Scan(&product.ProductID, &product.Name, &product.Price, &product.CreatedBy, &product.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, rest_err.NewBadRequestError(fmt.Sprintf("Product dengan product_id %d tidak ditemukan", input.ProductID))
		} else {
			return nil, rest_err.NewInternalServerError("gagal mengedit product", err)
		}
	}
	return &product, nil
}

func (u *productDao) Delete(productID int64) rest_err.APIError {
	sqlStatement := `
	DELETE FROM products 
	WHERE product_id = $1;
	`
	res, err := db.DB.Exec(context.Background(), sqlStatement, productID)
	if err != nil {
		return rest_err.NewInternalServerError("gagal saat penghapusan product", err)
	}
	if res.RowsAffected() != 1 {
		return rest_err.NewBadRequestError(fmt.Sprintf("Product dengan product_id %d tidak ditemukan", productID))
	}

	return nil
}

func (u *productDao) Get(productID int64) (*dto.Product, rest_err.APIError) {

	sqlStatement := `
	SELECT product_id, name, price,created_by, created_at 
	FROM products 
	WHERE product_id = $1;
	`
	row := db.DB.QueryRow(context.Background(), sqlStatement, productID)

	var product dto.Product
	err := row.Scan(&product.ProductID, &product.Name, &product.Name, &product.Price, &product.CreatedBy, &product.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, rest_err.NewBadRequestError(fmt.Sprintf("Product dengan ID %d tidak ditemukan", productID))
		} else {
			return nil, rest_err.NewInternalServerError("gagal mendapatkan detil product", err)
		}
	}
	return &product, nil
}

func (u *productDao) Find() ([]dto.Product, rest_err.APIError) {
	rows, err := db.DB.Query(context.Background(),
		"SELECT product_id,  name, price, created_by, created_at FROM products;")
	if err != nil {
		return nil, rest_err.NewInternalServerError("gagal mendapatkan daftar product", err)
	}

	var products []dto.Product
	for rows.Next() {
		product := dto.Product{}
		err := rows.Scan(&product.ProductID, &product.Name, &product.Price, &product.CreatedBy, &product.CreatedAt)
		if err != nil {
			return nil, rest_err.NewInternalServerError("gagal scan list product", err)
		}
		products = append(products, product)
	}
	defer rows.Close()
	return products, nil
}
