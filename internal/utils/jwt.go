package utils

import (
    "time"
    "strconv"

    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID int    `json:"userId"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func GenerateAccessToken(userID int, email, role, secret string) (string, error) {
    claims := &Claims{
        UserID: userID,
        Email:  email,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Minute)), // 15 минут
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "paydeya-backend",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(userID int, secret string) (string, error) {
    claims := &jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 дней
        IssuedAt:  jwt.NewNumericDate(time.Now()),
        Issuer:    "paydeya-backend",
        Subject:   strconv.Itoa(userID),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString, secret string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, jwt.ErrSignatureInvalid
    }

    return claims, nil
}