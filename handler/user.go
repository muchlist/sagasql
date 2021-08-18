package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/sagasql/dto"
	"github.com/muchlist/sagasql/service"
	"github.com/muchlist/sagasql/utils/mjwt"
	"github.com/muchlist/sagasql/utils/rest_err"
	"time"
)

func NewUserHandler(userService service.UserServiceAssumer) *userHandler {
	return &userHandler{
		service: userService,
	}
}

type userHandler struct {
	service service.UserServiceAssumer
}

// Login login
func (u *userHandler) Login(c *fiber.Ctx) error {
	var login dto.UserLoginRequest
	if err := c.BodyParser(&login); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if login.Username == "" || login.Password == "" {
		apiErr := rest_err.NewBadRequestError("username atau password tidak boleh kosong")
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	response, apiErr := u.service.Login(login)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": response})
}

// Register menambahkan user
func (u *userHandler) Register(c *fiber.Ctx) error {
	var user dto.UserRegisterReq
	if err := c.BodyParser(&user); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := user.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	insertUsername, apiErr := u.service.InsertUser(dto.User{
		Username:  dto.UppercaseString(user.Username),
		Email:     user.Email,
		Name:      user.Name,
		Password:  user.Password,
		Role:      user.Role,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	})
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	res := fmt.Sprintf("Register berhasil, ID: %s", *insertUsername)
	return c.JSON(fiber.Map{"error": nil, "data": res})
}

// Edit mengedit user
func (u *userHandler) Edit(c *fiber.Ctx) error {
	username := c.Params("username")

	var user dto.User
	user.Username = dto.UppercaseString(username)
	if err := c.BodyParser(&user); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	userEdited, apiErr := u.service.EditUser(user)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": userEdited})
}

// RefreshToken
func (u *userHandler) RefreshToken(c *fiber.Ctx) error {
	var payload dto.UserRefreshTokenRequest
	if err := c.BodyParser(&payload); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	response, apiErr := u.service.Refresh(payload)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": response})
}

// Delete menghapus user, idealnya melalui middleware is_admin
func (u *userHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	username := c.Params("username")

	if claims.Identity == username {
		apiErr := rest_err.NewBadRequestError("Tidak dapat menghapus akun terkait (diri sendiri)!")
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	apiErr := u.service.DeleteUser(username)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("user %s berhasil dihapus", username)})
}

// Get menampilkan user berdasarkan username
func (u *userHandler) Get(c *fiber.Ctx) error {
	userName := c.Params("username")
	user, apiErr := u.service.GetUser(userName)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": user})
}

// GetProfile mengembalikan user yang sedang login
func (u *userHandler) GetProfile(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	user, apiErr := u.service.GetUser(claims.Identity)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": user})
}

// Find menampilkan list user
func (u *userHandler) Find(c *fiber.Ctx) error {
	userList, apiErr := u.service.FindUsers()
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if userList == nil {
		userList = []dto.User{}
	}
	return c.JSON(fiber.Map{"error": nil, "data": userList})
}
