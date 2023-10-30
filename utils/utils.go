package utils

import (
	"crypto/rand"
	"fmt"
	"loa/user_content/types"
	"math/big"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)
 var Timeout = 60*3
func GetEnv() *types.Env{
  err := godotenv.Load()
  if err != nil{
	panic(err)
  }
  env := &types.Env{
     User: os.Getenv("DB_USER"),
	 Password: os.Getenv("DB_PASSWORD"),
	 DB: os.Getenv("DB_NAME"),
	 PORT: os.Getenv("PORT"),
	 URL: os.Getenv("URL"),
  }
 fmt.Println(env)
 return env
}
func HashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
    if err != nil{
		return "", err
	}
	return string(hashPassword), nil
}
func MatchPassword(toCheck string, hashed string) bool{
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(toCheck))
	return err == nil
}
func CreateJWT(username string)(string,error){
    env:= GetEnv()
	claims:= &jwt.MapClaims{
		"username": username,
		"expiresAt": time.Minute*3,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	key :=env.JwtSecret
	return token.SignedString([]byte(key))

}
func ValidateJWT(tokenString string)(*jwt.Token, error){
   env := GetEnv()
   secret:= env.JwtSecret
   return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
	if _,ok := t.Method.(*jwt.SigningMethodHMAC); !ok{
		return nil, fmt.Errorf("unexpected signing method: %v",t.Header["alg"])
	}
	return []byte(secret),nil
   })

}
func RanHash()string{
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var slug string
	for i:=0;i<8;i++{
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(characters))))
		if err != nil{
			fmt.Println("error generating random number")
			return ""
		}
		slug += string(characters[randomIndex.Int64()])
	}
	return slug
}