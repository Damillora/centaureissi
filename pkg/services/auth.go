package services

import (
	"errors"
	"time"

	"github.com/Damillora/centaureissi/pkg/config"
	"github.com/Damillora/centaureissi/pkg/database/schema"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func (cs *CentaureissiService) Login(username string, password string) *schema.User {
	user, err := cs.repository.GetUserByUsername(username)
	if err != nil {
		return nil
	}
	if user == nil {
		return nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil
	}
	return user
}

func (cs *CentaureissiService) CreateToken(user *schema.User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": user.Username,
		"iss":  "centaureissi-api",
		"sub":  user.ID,
		"aud":  "centaureissi",
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})
	jwtToken, _ := token.SignedString([]byte(config.CurrentConfig.AuthSecret))
	return jwtToken
}

func (cs *CentaureissiService) ValidateToken(signedToken string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(
		signedToken,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(config.CurrentConfig.AuthSecret), nil
		},
	)
	if err != nil {
		return nil, err
	}

	if !claims.VerifyExpiresAt(time.Now().Local().Unix(), true) {
		return nil, errors.New("token is expired")
	}

	return claims, nil
}
