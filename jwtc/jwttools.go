package jwtc

import (
	"context"
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
)

type TokenInfo struct {
	Uid      string `json:"uid"`
	SchoolID int64  `json:"school_id"`
	RoleType int64  `json:"role_type"`
}

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
