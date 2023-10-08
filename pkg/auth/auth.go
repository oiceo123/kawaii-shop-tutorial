package auth

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oiceo123/kawaii-shop-tutorial/config"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/users"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
	Admin   TokenType = "admin"
	ApiKey  TokenType = "apikey"
)

type kawaiiAuth struct {
	mapClaims *kawaiiMapClaims // mapClaims = payload
	cfg       config.IJwtConfig
}

type kawaiiAdmin struct {
	*kawaiiAuth
}

type kawaiiMapClaims struct {
	Claims *users.UserClaims `json:"claims"`
	jwt.RegisteredClaims
}

type IkawaiiAuth interface {
	SignToken() string
}

type IkawaiiAdmin interface {
	SignToken() string
}

func jwtTimeDuration(t int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func (a *kawaiiAuth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.SecretKey())
	return ss
}

func (a *kawaiiAdmin) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.AdminKey())
	return ss
}

func ParseToken(cfg config.IJwtConfig, tokenString string) (*kawaiiMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &kawaiiMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		// วิธีแปลงจาก type any เป็น type ตามที่เราต้องการให้เติม .(type ที่ต้องการจะเปลี่ยน)
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.SecretKey(), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token had expired")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	// วิธีแปลงจาก type any เป็น type ตามที่เราต้องการให้เติม .(type ที่ต้องการจะเปลี่ยน)
	if claims, ok := token.Claims.(*kawaiiMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}

func ParseAdminToken(cfg config.IJwtConfig, tokenString string) (*kawaiiMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &kawaiiMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		// วิธีแปลงจาก type any เป็น type ตามที่เราต้องการให้เติม .(type ที่ต้องการจะเปลี่ยน)
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.AdminKey(), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token had expired")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	// วิธีแปลงจาก type any เป็น type ตามที่เราต้องการให้เติม .(type ที่ต้องการจะเปลี่ยน)
	if claims, ok := token.Claims.(*kawaiiMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}

func RepeatToken(cfg config.IJwtConfig, claims *users.UserClaims, exp int64) string {
	obj := &kawaiiAuth{
		cfg: cfg,
		mapClaims: &kawaiiMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "kawaiishop-api",               // เว็บหรือบริษัทเจ้าของ token
				Subject:   "refresh-token",                // subject ของ token
				Audience:  []string{"customer", "admin"},  // ผู้รับ token
				ExpiresAt: jwtTimeRepeatAdapter(exp),      // เวลาหมดอายุของ token
				NotBefore: jwt.NewNumericDate(time.Now()), // เป็นเวลาที่บอกว่า token จะเริ่มใช้งานได้เมื่อไหร่
				IssuedAt:  jwt.NewNumericDate(time.Now()), // ใช้เก็บเวลาที่ token นี้เกิดปัญหา
			},
		},
	}
	return obj.SignToken()
}

func NewKawaiiAuth(tokenType TokenType, cfg config.IJwtConfig, claims *users.UserClaims) (IkawaiiAuth, error) {
	switch tokenType {
	case Access:
		return newAccessToken(cfg, claims), nil
	case Refresh:
		return newRefreshToken(cfg, claims), nil
	case Admin:
		return newAdminToken(cfg), nil
	default:
		return nil, fmt.Errorf("unknown token type")
	}
}

func newAccessToken(cfg config.IJwtConfig, claims *users.UserClaims) IkawaiiAuth {
	return &kawaiiAuth{
		cfg: cfg,
		mapClaims: &kawaiiMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "kawaiishop-api",                       // เว็บหรือบริษัทเจ้าของ token
				Subject:   "access-token",                         // subject ของ token
				Audience:  []string{"customer", "admin"},          // ผู้รับ token
				ExpiresAt: jwtTimeDuration(cfg.AccessExpiresAt()), // เวลาหมดอายุของ token
				NotBefore: jwt.NewNumericDate(time.Now()),         // เป็นเวลาที่บอกว่า token จะเริ่มใช้งานได้เมื่อไหร่
				IssuedAt:  jwt.NewNumericDate(time.Now()),         // ใช้เก็บเวลาที่ token นี้เกิดปัญหา
			},
		},
	}
}

func newRefreshToken(cfg config.IJwtConfig, claims *users.UserClaims) IkawaiiAuth {
	return &kawaiiAuth{
		cfg: cfg,
		mapClaims: &kawaiiMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "kawaiishop-api",                        // เว็บหรือบริษัทเจ้าของ token
				Subject:   "refresh-token",                         // subject ของ token
				Audience:  []string{"customer", "admin"},           // ผู้รับ token
				ExpiresAt: jwtTimeDuration(cfg.RefreshExpiresAt()), // เวลาหมดอายุของ token
				NotBefore: jwt.NewNumericDate(time.Now()),          // เป็นเวลาที่บอกว่า token จะเริ่มใช้งานได้เมื่อไหร่
				IssuedAt:  jwt.NewNumericDate(time.Now()),          // ใช้เก็บเวลาที่ token นี้เกิดปัญหา
			},
		},
	}
}

func newAdminToken(cfg config.IJwtConfig) IkawaiiAuth {
	return &kawaiiAdmin{
		kawaiiAuth: &kawaiiAuth{
			cfg: cfg,
			mapClaims: &kawaiiMapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "kawaiishop-api",               // เว็บหรือบริษัทเจ้าของ token
					Subject:   "admin-token",                  // subject ของ token
					Audience:  []string{"admin"},              // ผู้รับ token
					ExpiresAt: jwtTimeDuration(300),           // เวลาหมดอายุของ token
					NotBefore: jwt.NewNumericDate(time.Now()), // เป็นเวลาที่บอกว่า token จะเริ่มใช้งานได้เมื่อไหร่
					IssuedAt:  jwt.NewNumericDate(time.Now()), // ใช้เก็บเวลาที่ token นี้เกิดปัญหา
				},
			},
		},
	}
}
