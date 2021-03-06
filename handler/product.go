package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/sagasql/dto"
	"github.com/muchlist/sagasql/service"
	"github.com/muchlist/sagasql/utils/mjwt"
	"github.com/muchlist/sagasql/utils/rest_err"
	"strconv"
	"time"
)

func NewProductHandler(productService service.ProductServiceAssumer) *productHandler {
	return &productHandler{
		service: productService,
	}
}

type productHandler struct {
	service service.ProductServiceAssumer
}

// Insert menambahkan product
func (u *productHandler) Insert(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var product dto.ProductReq
	if err := c.BodyParser(&product); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := product.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	insertProductID, apiErr := u.service.InsertProduct(dto.Product{
		Name:      dto.UppercaseString(product.Name),
		Price:     product.Price,
		CreatedBy: claims.Identity,
		CreatedAt: time.Now().Unix(),
	})
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	res := fmt.Sprintf("Register berhasil, ID: %d", *insertProductID)
	return c.JSON(fiber.Map{"error": nil, "data": res})
}

// Edit mengedit product
func (u *productHandler) Edit(c *fiber.Ctx) error {
	productIDStr := c.Params("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		apiErr := rest_err.NewBadRequestError("ID harus dalam bentuk angka")
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	var product dto.Product
	product.ProductID = productID
	if err := c.BodyParser(&product); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	productEdited, apiErr := u.service.EditProduct(product)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": productEdited})
}

// Delete menghapus product, idealnya melalui middleware is_admin
func (u *productHandler) Delete(c *fiber.Ctx) error {
	productIDStr := c.Params("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		apiErr := rest_err.NewBadRequestError("ID harus dalam bentuk angka")
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	apiErr := u.service.DeleteProduct(productID)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("product %d berhasil dihapus", productID)})
}

// Get menampilkan product berdasarkan productID
func (u *productHandler) Get(c *fiber.Ctx) error {
	productIDStr := c.Params("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		apiErr := rest_err.NewBadRequestError("ID harus dalam bentuk angka")
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	product, apiErr := u.service.GetProduct(productID)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": product})
}

// Find menampilkan list product
func (u *productHandler) Find(c *fiber.Ctx) error {
	search := c.Query("search")

	productList, apiErr := u.service.FindProducts(search)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if productList == nil {
		productList = []dto.Product{}
	}
	return c.JSON(fiber.Map{"error": nil, "data": productList})
}

// UploadImage melakukan pengambilan file menggunakan form "image" mengecek ekstensi dan memasukkannya ke database
func (u *productHandler) UploadImage(c *fiber.Ctx) error {
	productIDStr := c.Params("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		apiErr := rest_err.NewBadRequestError("ID harus dalam bentuk angka")
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	randomName := fmt.Sprintf("%d-%d", productID, time.Now().Unix())
	// simpan image
	pathInDB, apiErr := saveImage(c, "product", randomName)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	// update path image di database
	productResult, apiErr := u.service.PutImage(productID, pathInDB)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": productResult})
}
