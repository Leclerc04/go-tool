package jwtc

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v4"
)

// TokenInfo 存储解析后的token信息
type TokenInfo struct {
	Uid      string `json:"uid"`
	SchoolID int64  `json:"school_id"`
	RoleType int64  `json:"role_type"`
}

// GenJwtToken 生成JWT token
func GenJwtToken(secretKey string, iat, expiration int64, roleType, userID string) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		//"schoolID": schoolID,
		"roleType": roleType,
		"iat":      iat,
		"exp":      iat + expiration,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// ParseJwtToken 从上下文中解析token信息
func ParseJwtToken(ctx context.Context) (tokenInfo TokenInfo, err error) {
	var (
		userId string
		//schoolIDJSONum json.Number
		roleTypeJSONum json.Number
		//schoolID       int64
		roleType int64
		ok       bool
	)
	//获取uid
	if userId, ok = ctx.Value("userID").(string); !ok {
		return
	}
	//获取schoolID
	//if schoolIDJSONum, ok = ctx.Value("schoolID").(json.Number); !ok {
	//	return
	//}

	//if schoolID, err = schoolIDJSONum.Int64(); err != nil {
	//	return
	//}

	// 获取roleType
	if roleTypeJSONum, ok = ctx.Value("roleType").(json.Number); !ok {
		return
	}

	if roleType, err = roleTypeJSONum.Int64(); err != nil {
		return
	}
	tokenInfo = TokenInfo{
		Uid: userId,
		//SchoolID: schoolID,
		RoleType: roleType,
	}
	return
}

// ParseTokenToUID 解析token字符串获取uid
func ParseTokenToUID(tokenString, secretKey string) (string, error) {
	// 解析token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", err
	}

	// 提取claims中的userID
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["userID"].(string); ok {
			return userID, nil
		}
		return "", errors.New("userID not found in token")
	}

	return "", errors.New("invalid token")
}
