package dao

import (
	"context"
	"fmt"
	"github.com/muchlist/sagasql/db"
	"github.com/muchlist/sagasql/dto"
	"github.com/muchlist/sagasql/utils/rest_err"
	"github.com/muchlist/sagasql/utils/sql_err"
)

func NewProductDao() ProductDaoAssumer {
	return &productDao{}
}

type ProductDaoAssumer interface {
	Insert(product dto.Product) (*int64, rest_err.APIError)
	Edit(productInput dto.Product) (*dto.Product, rest_err.APIError)
	Delete(productID int64) rest_err.APIError
	UploadImage(productID int64, imagePath string) (*dto.Product, rest_err.APIError)
	Get(productID int64) (*dto.Product, rest_err.APIError)
	Find() ([]dto.Product, rest_err.APIError)
	Search(productName dto.UppercaseString) ([]dto.Product, rest_err.APIError)
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
		return nil, sql_err.ParseError(err)
	}
	return &productID, nil
}

func (u *productDao) Edit(input dto.Product) (*dto.Product, rest_err.APIError) {
	sqlStatement := `
	UPDATE products 
	SET name = $2, price = $3
	WHERE product_id = $1 
	RETURNING product_id, name, price, image, created_by, created_at;
	`

	var product dto.Product
	err := db.DB.QueryRow(
		context.Background(),
		sqlStatement, input.ProductID, input.Name, input.Price,
	).Scan(&product.ProductID, &product.Name, &product.Price, &product.Image, &product.CreatedBy, &product.CreatedAt)
	if err != nil {
		return nil, sql_err.ParseError(err)
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
		return sql_err.ParseError(err)
	}
	if res.RowsAffected() != 1 {
		return rest_err.NewBadRequestError(fmt.Sprintf("Product dengan product_id %d tidak ditemukan", productID))
	}

	return nil
}

func (u *productDao) UploadImage(productID int64, imagePath string) (*dto.Product, rest_err.APIError) {
	sqlStatement := `
	UPDATE products 
	SET image = $2
	WHERE product_id = $1 
	RETURNING product_id, name, price, image, created_by, created_at;
	`

	var product dto.Product
	err := db.DB.QueryRow(
		context.Background(),
		sqlStatement, productID, imagePath,
	).Scan(&product.ProductID, &product.Name, &product.Price, &product.Image, &product.CreatedBy, &product.CreatedAt)
	if err != nil {
		return nil, sql_err.ParseError(err)
	}
	return &product, nil
}

func (u *productDao) Get(productID int64) (*dto.Product, rest_err.APIError) {

	sqlStatement := `
	SELECT product_id, name, price, image, created_by, created_at 
	FROM products 
	WHERE product_id = $1;
	`
	row := db.DB.QueryRow(context.Background(), sqlStatement, productID)

	var product dto.Product
	err := row.Scan(&product.ProductID, &product.Name, &product.Price, &product.Image, &product.CreatedBy, &product.CreatedAt)
	if err != nil {
		return nil, sql_err.ParseError(err)
	}
	return &product, nil
}

func (u *productDao) Find() ([]dto.Product, rest_err.APIError) {
	rows, err := db.DB.Query(context.Background(),
		`	SELECT product_id, name, price, image, created_by, created_at 
				FROM products 
				ORDER BY name ASC;`)
	if err != nil {
		return nil, rest_err.NewInternalServerError("gagal mendapatkan daftar product", err)
	}
	defer rows.Close()
	var products []dto.Product
	for rows.Next() {
		product := dto.Product{}
		err := rows.Scan(&product.ProductID, &product.Name, &product.Price, &product.Image, &product.CreatedBy, &product.CreatedAt)
		if err != nil {
			return nil, sql_err.ParseError(err)
		}
		products = append(products, product)
	}
	return products, nil
}

func (u *productDao) Search(productName dto.UppercaseString) ([]dto.Product, rest_err.APIError) {

	rows, err := db.DB.Query(context.Background(),
		`SELECT product_id, name, price, image, created_by, created_at FROM products WHERE name LIKE '%'|| $1 || '%' ORDER BY name ASC ;`, productName)
	if err != nil {
		return nil, rest_err.NewInternalServerError("gagal mendapatkan daftar product", err)
	}

	defer rows.Close()

	var products []dto.Product
	for rows.Next() {
		product := dto.Product{}
		err := rows.Scan(&product.ProductID, &product.Name, &product.Price, &product.Image, &product.CreatedBy, &product.CreatedAt)
		if err != nil {
			return nil, sql_err.ParseError(err)
		}
		products = append(products, product)
	}
	return products, nil
}
