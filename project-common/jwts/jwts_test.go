package jwts

import (
	"fmt"
	"strconv"
	"testing"
)

func TestParseToken(t *testing.T) {
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTA4MjMzMjUsInRva2VuIjoiMjgifQ.g5h6fs71mwjeIFQWHebnW1eCeIge-mus6o_8Rz64swE"
	//config.C.JwtConfig.AccessSecret
	parseToken, err := ParseToken(jwt, "msproject") //返回的val实际上是空的
	if err != nil {
		fmt.Println(err)
	}
	id, _ := strconv.ParseInt(parseToken, 10, 64)
	fmt.Println("parse解析的结果为:", id)
}
