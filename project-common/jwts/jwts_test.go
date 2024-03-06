package jwts

import (
	"fmt"
	"testing"
)

func TestParseToken(t *testing.T) {
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTAwNTg0OTcsInRva2VuIjoiXHUwMDEzIn0.WshbQF1DuWPiWbrL8tDUqhDlpnJzbyDYXxYdP6OSlig"
	//config.C.JwtConfig.AccessSecret
	claims := parseToken(jwt, "msproject")
	fmt.Println(claims["id"])

}
