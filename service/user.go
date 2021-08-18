package user_service

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
	GetUser(userID int64) (*dto.User, rest_err.APIError)
	FindUsers() ([]dto.User, rest_err.APIError)
	InsertUser(user dto.User) (*int64, rest_err.APIError)
	EditUser(request dto.User) (*dto.User, rest_err.APIError)
	DeleteUser(userID int64) rest_err.APIError
	Login(login dto.User) (*dto.UserLoginResponse, rest_err.APIError)
	Refresh(payload dto.UserRefreshTokenRequest) (*dto.UserRefreshTokenResponse, rest_err.APIError)
}

// GetUser mendapatkan user dari database
func (u *userService) GetUser(userID int64) (*dto.User, rest_err.APIError) {
	user, err := u.dao.Get(userID)
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

// InsertUser melakukan register user
func (u *userService) InsertUser(user dto.User) (*int64, rest_err.APIError) {
	hashPassword, err := u.crypto.GenerateHash(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashPassword
	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = time.Now().Unix()

	insertedID, err := u.dao.Insert(user)
	if err != nil {
		return nil, err
	}
	return insertedID, nil
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

// DeleteUser
func (u *userService) DeleteUser(userID int64) rest_err.APIError {
	err := u.dao.Delete(userID)
	if err != nil {
		return err
	}
	return nil
}

// Login
func (u *userService) Login(login dto.User) (*dto.UserLoginResponse, rest_err.APIError) {
	user, err := u.dao.Get(login.UserID)
	if err != nil {
		return nil, err
	}

	if !u.crypto.IsPWAndHashPWMatch(login.Password, user.Password) {
		return nil, rest_err.NewUnauthorizedError("Username atau password tidak valid")
	}

	AccessClaims := mjwt.CustomClaim{
		Identity:    user.UserID,
		Name:        user.Name,
		UserName:    user.Username,
		Roles:       user.Role,
		ExtraMinute: 60 * 24 * 1, // 1 Hour
		Type:        mjwt.Access,
		Fresh:       true,
	}

	RefreshClaims := mjwt.CustomClaim{
		Identity:    user.UserID,
		Name:        user.Name,
		UserName:    user.Username,
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
		UserID:       user.UserID,
		Username:     user.Username,
		Email:        user.Email,
		Name:         user.Name,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expired:      time.Now().Add(time.Minute * time.Duration(60*24*1)).Unix(),
	}

	return &userResponse, nil
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
		Identity:    user.UserID,
		UserName:    user.Username,
		Name:        user.Username,
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
