/**
 *
 * @company: Co.预见（天津）智能科技有限公司
 * @Author:  ZhaoYi
 * @Date:    2023/2/20 11:45
 */

package duang

import (
	sf "github.com/bwmarrin/snowflake"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

var nodeX *sf.Node

// Init 雪花算法初始化
func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	sf.Epoch = st.UnixNano() / 1000000
	nodeX, err = sf.NewNode(machineID)
	return
}

// GenID 雪花算法返回
func GenID() string {
	return nodeX.Generate().String()
}

// CheckFileIsExist 检查目录是否存在
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, er := os.Stat(filename); os.IsNotExist(er) {
		exist = false
	}
	return exist
}

type MyCustomClaims struct {
	ID string
	jwt.RegisteredClaims
}

// EnToken 生成token
func EnToken(val string) (string, error) {
	jwtKey := []byte("boom")
	claims := MyCustomClaims{
		val,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
			Issuer:    "test",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, er := token.SignedString(jwtKey)
	if er != nil {
		return "", er
	}
	return tokenString, nil
}

// DeToken 解密token
func DeToken(tokenString string) (*MyCustomClaims, error) {
	jwtKey := []byte("boom")
	token, er := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if er != nil {
		return nil, er
	}
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, er
	}
}
