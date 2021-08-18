package mjwt

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/muchlist/sagasql/utils/rest_err"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	CLAIMS    = "claims"
	secretKey = "SECRET_KEY"

	identityKey  = "identity"
	nameKey      = "name"
	rolesKey     = "roles"
	tokenTypeKey = "type"
	expKey       = "exp"
	freshKey     = "fresh"
)

var (
	JwtObj JWTAssumer
	secret []byte
)

func NewJwt() JWTAssumer {
	return JwtObj
}

func init() {
	secret = []byte(os.Getenv(secretKey))
	if string(secret) == "" {
		log.Fatal("Secret key tidak boleh kosong, ENV : SECRET_KEY")
	}

	JwtObj = &jwtUtils{}
}

type JWTAssumer interface {
	GenerateToken(claims CustomClaim) (string, rest_err.APIError)
	ValidateToken(tokenString string) (*jwt.Token, rest_err.APIError)
	ReadToken(token *jwt.Token) (*CustomClaim, rest_err.APIError)
}

type jwtUtils struct {
}

// GenerateToken membuat token jwt untuk login header, untuk menguji nilai payloadnya
// dapat menggunakan situs jwt.io
func (j *jwtUtils) GenerateToken(claims CustomClaim) (string, rest_err.APIError) {
	expired := time.Now().Add(time.Minute * claims.ExtraMinute).Unix()

	jwtClaim := jwt.MapClaims{}
	jwtClaim[identityKey] = claims.Identity
	jwtClaim[nameKey] = claims.Name
	jwtClaim[rolesKey] = claims.Roles
	jwtClaim[expKey] = expired
	jwtClaim[tokenTypeKey] = claims.Type
	jwtClaim[freshKey] = claims.Fresh

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaim)

	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", rest_err.NewInternalServerError("gagal menandatangani token", err)
	}

	return signedToken, nil
}

// ReadToken membaca inputan token dan menghasilkan pointer struct CustomClaim
// struct CustomClaim digunakan untuk nilai passing antar middleware
func (j *jwtUtils) ReadToken(token *jwt.Token) (*CustomClaim, rest_err.APIError) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, rest_err.NewInternalServerError("gagal mapping token", nil)
	}

	customClaim := CustomClaim{
		Identity: claims[identityKey].(string),
		Name:     claims[nameKey].(string),
		Exp:      int64(claims[expKey].(float64)),
		Roles:    claims[rolesKey].(string),
		Type:     int(claims[tokenTypeKey].(float64)),
		Fresh:    claims[freshKey].(bool),
	}

	return &customClaim, nil
}

// ValidateToken memvalidasi apakah token string masukan valid, termasuk memvalidasi apabila field exp nya kadaluarsa
func (j *jwtUtils) ValidateToken(tokenString string) (*jwt.Token, rest_err.APIError) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, rest_err.NewAPIError("Token signing method salah", http.StatusUnprocessableEntity, "jwt_error", nil)
		}
		return secret, nil
	})

	// Jika expired akan muncul disini asalkan ada claims exp
	if err != nil {
		return nil, rest_err.NewAPIError("Token tidak valid", http.StatusUnprocessableEntity, "jwt_error", []interface{}{err.Error()})
	}

	return token, nil
}
