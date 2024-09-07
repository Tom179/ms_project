package jwts

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

type JwtToken struct {
	AccessToken  string
	RefreshToken string
	AccessExp    int64
	RefreshExp   int64
}

func CreateToken(val string, exp time.Duration, secret string, refreshExp time.Duration, refreshSecret string) *JwtToken {
	fmt.Println("生成token时value为:", val, "密钥为:", secret) ////////////////////////////
	aExp := time.Now().Add(exp).Unix()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token": val,
		"exp":   aExp,
	})
	aToken, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	rExp := time.Now().Add(refreshExp).Unix()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token": val,
		"exp":   rExp,
	})
	rToken, err := refreshToken.SignedString([]byte(refreshSecret))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &JwtToken{
		AccessExp:    aExp,
		AccessToken:  aToken,
		RefreshExp:   rExp,
		RefreshToken: rToken,
	}
}

func ParseToken(tokenString string, secret string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		val := claims["token"].(string)       //注意这个claim的类型和返回类型
		exp := int64(claims["exp"].(float64)) //?为什么不直接断言为int64?
		fmt.Println("val:", val, " exp:", exp)
		if exp <= time.Now().Unix() {
			return "", errors.New("token过期")
		}

		return val, nil
	} else {
		return "", err
	}

}
