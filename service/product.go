package service

import (
	"github.com/muchlist/sagasql/dao"
	"github.com/muchlist/sagasql/dto"
	"github.com/muchlist/sagasql/utils/rest_err"
)

func NewProductService(dao dao.ProductDaoAssumer) ProductServiceAssumer {
	return &productService{
		dao: dao,
	}
}

type productService struct {
	dao dao.ProductDaoAssumer
}

type ProductServiceAssumer interface {
	InsertProduct(product dto.Product) (*int64, rest_err.APIError)
	EditProduct(request dto.Product) (*dto.Product, rest_err.APIError)
	DeleteProduct(productID int64) rest_err.APIError
	GetProduct(productID int64) (*dto.Product, rest_err.APIError)
	FindProducts(search string) ([]dto.Product, rest_err.APIError)
}

// InsertProduct melakukan register product
func (u *productService) InsertProduct(product dto.Product) (*int64, rest_err.APIError) {
	insertedProductID, err := u.dao.Insert(product)
	if err != nil {
		return nil, err
	}
	return insertedProductID, nil
}

// EditProduct
func (u *productService) EditProduct(request dto.Product) (*dto.Product, rest_err.APIError) {
	result, err := u.dao.Edit(request)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteProduct
func (u *productService) DeleteProduct(productID int64) rest_err.APIError {
	err := u.dao.Delete(productID)
	if err != nil {
		return err
	}
	return nil
}

// GetProduct mendapatkan product dari database
func (u *productService) GetProduct(productID int64) (*dto.Product, rest_err.APIError) {
	product, err := u.dao.Get(productID)
	if err != nil {
		return nil, err
	}
	return product, nil
}

// FindProducts
func (u *productService) FindProducts(search string) ([]dto.Product, rest_err.APIError) {

	var productList []dto.Product
	var err rest_err.APIError
	if len(search) > 0 {
		productList, err = u.dao.Search(dto.UppercaseString(search))
	} else {
		productList, err = u.dao.Find()
	}
	if err != nil {
		return nil, err
	}
	return productList, nil
}
