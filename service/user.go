package service

import (
	"github.com/muchlist/sagasql/dao"
	"github.com/muchlist/sagasql/dto"
	"github.com/muchlist/sagasql/utils/mcrypt"
	"github.com/muchlist/sagasql/utils/mjwt"
	"github.com/muchlist/sagasql/utils/rest_err"
	"net/http"
	"time"
)

func NewUserService(dao dao.UserDaoAssumer, crypto mcrypt.BcryptAssumer, jwt mjwt.JWTAssumer) UserServiceAssumer {
	return &userService{
		dao:    dao,
		crypto: crypto,
		jwt:    jwt,
	}
}

type userService struct {
	dao    dao.UserDaoAssumer
	crypto mcrypt.BcryptAssumer
	jwt    mjwt.JWTAssumer
}

type UserServiceAssumer interface {
	Login(login dto.UserLoginRequest) (*dto.UserLoginResponse, rest_err.APIError)
	InsertUser(user dto.User) (*string, rest_err.APIError)
	EditUser(request dto.User) (*dto.User, rest_err.APIError)
	Refresh(payload dto.UserRefreshTokenRequest) (*dto.UserRefreshTokenResponse, rest_err.APIError)
	DeleteUser(username string) rest_err.APIError
	GetUser(username string) (*dto.User, rest_err.APIError)
	FindUsers() ([]dto.User, rest_err.APIError)
}

// Login
func (u *userService) Login(login dto.UserLoginRequest) (*dto.UserLoginResponse, rest_err.APIError) {
	user, err := u.dao.Get(login.Username)
	if err != nil {
		return nil, rest_err.NewBadRequestError("Username atau password tidak valid")
	}

	if !u.crypto.IsPWAndHashPWMatch(login.Password, user.Password) {
		return nil, rest_err.NewUnauthorizedError("Username atau password tidak valid")
	}

	AccessClaims := mjwt.CustomClaim{
		Identity:    string(user.Username),
		Name:        user.Name,
		Roles:       user.Role,
		ExtraMinute: 60 * 24 * 1, // 1 Hour
		Type:        mjwt.Access,
		Fresh:       true,
	}

	RefreshClaims := mjwt.CustomClaim{
		Identity:    string(user.Username),
		Name:        user.Name,
		Roles:       user.Role,
		ExtraMinute: 60 * 24 * 10, // 60 days
		Type:        mjwt.Refresh,
	}

	accessToken, err := u.jwt.GenerateToken(AccessClaims)
	if err != nil {
		return nil, err
	}
	refreshToken, err := u.jwt.GenerateToken(RefreshClaims)
	if err != nil {
		return nil, err
	}

	userResponse := dto.UserLoginResponse{
		Username:     string(user.Username),
		Email:        user.Email,
		Name:         user.Name,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expired:      time.Now().Add(time.Minute * time.Duration(60*24*1)).Unix(),
	}

	return &userResponse, nil
}

// InsertUser melakukan register user
func (u *userService) InsertUser(user dto.User) (*string, rest_err.APIError) {
	hashPassword, err := u.crypto.GenerateHash(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashPassword
	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = time.Now().Unix()

	insertedUserID, err := u.dao.Insert(user)
	if err != nil {
		return nil, err
	}
	return insertedUserID, nil
}

// EditUser
func (u *userService) EditUser(request dto.User) (*dto.User, rest_err.APIError) {
	request.UpdatedAt = time.Now().Unix()
	result, err := u.dao.Edit(request)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Refresh token
func (u *userService) Refresh(payload dto.UserRefreshTokenRequest) (*dto.UserRefreshTokenResponse, rest_err.APIError) {
	token, apiErr := u.jwt.ValidateToken(payload.RefreshToken)
	if apiErr != nil {
		return nil, apiErr
	}
	claims, apiErr := u.jwt.ReadToken(token)
	if apiErr != nil {
		return nil, apiErr
	}

	// cek apakah tipe claims token yang dikirim adalah tipe refresh (1)
	if claims.Type != mjwt.Refresh {
		return nil, rest_err.NewAPIError("Token tidak valid", http.StatusUnprocessableEntity, "jwt_error", []interface{}{"not a refresh token"})
	}

	// mendapatkan data terbaru dari user
	user, apiErr := u.dao.Get(claims.Identity)
	if apiErr != nil {
		return nil, apiErr
	}

	AccessClaims := mjwt.CustomClaim{
		Identity:    string(user.Username),
		Name:        user.Name,
		Roles:       user.Role,
		ExtraMinute: time.Duration(60 * 60 * 1),
		Type:        mjwt.Access,
		Fresh:       false,
	}

	accessToken, err := u.jwt.GenerateToken(AccessClaims)
	if err != nil {
		return nil, err
	}

	userRefreshTokenResponse := dto.UserRefreshTokenResponse{
		AccessToken: accessToken,
		Expired:     time.Now().Add(time.Minute * time.Duration(60*60*1)).Unix(),
	}

	return &userRefreshTokenResponse, nil
}

// DeleteUser
func (u *userService) DeleteUser(userName string) rest_err.APIError {
	err := u.dao.Delete(userName)
	if err != nil {
		return err
	}
	return nil
}

// GetUser mendapatkan user dari database
func (u *userService) GetUser(userName string) (*dto.User, rest_err.APIError) {
	user, err := u.dao.Get(userName)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindUsers
func (u *userService) FindUsers() ([]dto.User, rest_err.APIError) {
	userList, err := u.dao.Find()
	if err != nil {
		return nil, err
	}
	return userList, nil
}
