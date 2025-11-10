package jwtauth

import (
	"errors"
	"fmt"

	"booking.com/internal/db/postgresql/dao"
	"booking.com/internal/db/postgresql/model"
	"booking.com/pkg/constants"
	"booking.com/pkg/utils"
	"github.com/golang-jwt/jwt"
)

// GetToken generates a signed JWT for a given username
func GetToken(userName, salt string, exp int64) (string, error) {
	claims := jwt.MapClaims{
		constants.UserName: userName,
		constants.Expiry:   utils.GetExpTime(exp),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(salt))
	if err != nil {
		return "", err
	}
	return signed, nil
}

// IsTokenValid validates a JWT and checks user + expiry
func IsTokenValid(token string, refresh bool) (bool, *model.User, error) {
	tempToken, _, _ := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
	claims := getJwtClaims(tempToken)
	if claims == nil {
		return false, nil, errors.New("invalid token structure")
	}

	userName, ok := claims[constants.UserName].(string)
	if !ok {
		return false, nil, errors.New("username missing in claims")
	}

	q := dao.User
	user, _ := q.Where(q.Email.Eq(userName)).Or(q.Phone.Eq(userName)).First()
	if user == nil {
		return false, nil, errors.New("invalid user")
	}
	jwtToken, err := getJwtTokenData(token, user.Salt)
	if err != nil || !jwtToken.Valid {
		return false, nil, fmt.Errorf("invalid token: %v", err)
	}
	if !refresh {
		claims := getJwtClaims(jwtToken)
		if claims == nil {
			return false, nil, errors.New("invalid token claims")
		}
		exp, ok := claims[constants.Expiry].(float64)
		if !ok {
			return false, nil, errors.New("invalid expiry claim type")
		}
		if utils.IsExpired(int64(exp)) {
			return false, nil, errors.New("token expired")
		}
	}

	return true, user, nil
}

func RefreshToken(token string, refresh bool) (string, error) {
	valid, user, err := IsTokenValid(token, refresh)
	if err != nil || !valid {
		return "", fmt.Errorf("invalid token, error: %v", err)
	}
	return GetToken(user.Email, user.Salt, 15)
}

func getJwtTokenData(token, salt string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", t.Header["alg"])
		}
		return []byte(salt), nil
	})
}
func getJwtClaims(token *jwt.Token) jwt.MapClaims {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil
	}
	return claims
}
