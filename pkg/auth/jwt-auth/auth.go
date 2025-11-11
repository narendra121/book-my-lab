package jwtauth

import (
	"errors"
	"fmt"

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
func IsTokenValid(token, salt string, validate func(...string) bool) (bool, error) {
	jwtToken, err := getJwtTokenData(token, salt)
	if err != nil || !jwtToken.Valid {
		return false, fmt.Errorf("invalid token: %v", err)
	}
	claims := getJwtClaims(jwtToken)
	if claims == nil {
		return false, errors.New("invalid token claims")
	}
	userName, ok := claims[constants.UserName].(string)
	if !ok {
		return false, errors.New("username missing in claims")
	}
	if validate != nil {
		if !validate(userName) {
			return false, errors.New("failed custom validation of tokenr")
		}
	}
	exp, ok := claims[constants.Expiry].(float64)
	if !ok {
		return false, errors.New("invalid expiry claim type")
	}
	if utils.IsExpired(int64(exp)) {
		return false, errors.New("token expired")
	}
	return true, nil
}
func GetUnVerifiedJwtClaims(token, claimKey string) (string, error) {
	jwtToken, _, _ := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
	if jwtToken == nil {
		return "", fmt.Errorf("invalid token structure")
	}
	claims := getJwtClaims(jwtToken)
	if claims == nil {
		return "", fmt.Errorf("invalid token structure")
	}
	val, ok := claims[claimKey]
	if !ok || val == "" {
		return "", fmt.Errorf("claim not found")
	}
	return val.(string), nil
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
