package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/sagasql/utils/rest_err"
	"net/http"
	"path/filepath"
	"strings"
)

const (
	jpgExtension  = ".jpg"
	pngExtension  = ".png"
	jpegExtension = ".jpeg"
)

// saveImage return path to save in db
func saveImage(c *fiber.Ctx, folder string, imageName string) (string, rest_err.APIError) {
	file, err := c.FormFile("image")
	if err != nil {
		apiErr := rest_err.NewAPIError("File gagal di upload", http.StatusBadRequest, "bad_request", []interface{}{err.Error()})
		return "", apiErr
	}

	fileName := file.Filename
	fileExtension := strings.ToLower(filepath.Ext(fileName))
	if !(fileExtension == jpgExtension || fileExtension == pngExtension || fileExtension == jpegExtension) {
		apiErr := rest_err.NewBadRequestError("Ektensi file tidak di support")
		return "", apiErr
	}

	if file.Size > 2*1024*1024 { // 1 MB
		apiErr := rest_err.NewBadRequestError("Ukuran file tidak dapat melebihi 2MB")
		return "", apiErr
	}

	// rename image
	// path := filepath.Join("static", "image", folder, imageName + fileExtension)
	// pathInDB := filepath.Join("image", folder, imageName + fileExtension)
	path := fmt.Sprintf("static/image/%s/%s", folder, imageName+fileExtension)
	pathInDB := fmt.Sprintf("image/%s/%s", folder, imageName+fileExtension)

	err = c.SaveFile(file, path)
	if err != nil {
		apiErr := rest_err.NewInternalServerError("File gagal di upload", err)
		return "", apiErr
	}

	return pathInDB, nil
}
