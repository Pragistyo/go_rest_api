package controllers

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"encoding/json"
	"strings"
	"context"

	jwt "github.com/dgrijalva/jwt-go"
)

type M map[string]interface{}

func MiddlewareJWTAuthorization(next http.HandlerFunc) http.HandlerFunc  {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	  // Our middleware logic goes here...

	//   if r.URL.Path != "/verifyToken/" {
	// 	next.ServeHTTP(w, r)
	// 	return
	// }

	  var tokenString string

	  authorizationHeader := r.Header.Get("Authorization")
  
	  if !strings.Contains(authorizationHeader, "Bearer") {
		  // http.Error(w, "Invalid token", http.StatusBadRequest)
		  w.WriteHeader(http.StatusForbidden)
	  }
  
	  tokenString = strings.Replace(authorizationHeader, "Bearer: ", "", -1)
	  fmt.Println("====== middleware ======")
	  tokenInfo, err :=VerifyToken(w,r,tokenString) //tokenInfo to context

	  if err!= nil{
		  w.WriteHeader(http.StatusUnauthorized)
		  respString, _ := json.Marshal( M{ "Message": "invalid token", "Status": 401 } )
		  w.Write([]byte  (respString))
		  return
	  }

	  ctx := context.WithValue(context.Background(), "tokenData", tokenInfo)
	  r = r.WithContext(ctx)

	  next.ServeHTTP(w, r)
	})
  }

type tokenStruct struct {
	Id 			float64 		`json:"id"`
	Username	string		`json:"username"`
	Authority 	float64		`json:"authority"`
}

func authToken( tokenString string ) (*jwt.Token, error) {
	fmt.Println(" ==== masuk auth token ===== ")

	

	token,err := jwt.Parse( tokenString, func( token *jwt.Token) (interface{}, error) {
		fmt.Println(" ==== masuk auth token 1,5 ===== ")
		fmt.Println(tokenString)

		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			fmt.Println(" ==== masuk auth token 2 ===== ")

			return nil, fmt.Errorf("Signing method invalid")

		} else if method != jwt.SigningMethodHS256 {

			fmt.Println(" ==== masuk auth token 3 ===== ")

			return nil, fmt.Errorf("Signing method invalid")

		}
		fmt.Println(" ==== masuk auth token 4 ===== ")
		secretKey := []byte( os.Getenv( "SECRET_KEY_JWT" ) )

		fmt.Println( reflect.TypeOf(secretKey) )
		return secretKey, nil
		
	})

	if err!= nil {
		fmt.Println(" ==== masuk auth token ===== 10: ", err)
			return nil, err
	}
	fmt.Println(" ===== masuk auth token 11 ===== ")
	

	return token, nil
}


func validToken(tokenString string, token *jwt.Token) bool{
	fmt.Println("====== masuk valid Token ======")
	token, err := authToken(tokenString)
	if err != nil {
		fmt.Println("====== masuk valid Token 2 ======")
		return false
	}

	_, ok := token.Claims.(jwt.MapClaims)
	 if !ok || !token.Valid {
		fmt.Println("====== masuk valid Token 3 ======")
		return false
	 }
	 fmt.Println("====== masuk valid Token 4 ======")
	return true
}

  
func VerifyToken(w http.ResponseWriter,r *http.Request,tokenString string)( jwt.MapClaims, error) {

	
	token, err := authToken(tokenString)
	if err != nil {
		// w.WriteHeader(http.StatusUnauthorized)
		return nil,err
	}

	if !validToken( tokenString, token) {
		// w.WriteHeader(http.StatusUnauthorized)
		return nil, err
	}

	fmt.Println(" ===== VerifyToken  ===== ")
	fmt.Println(token)
	tokenInfo := token.Claims.(jwt.MapClaims)
	fmt.Println(" ===== token claims ===== ")
	
	auth, _ := tokenInfo["authority"]

	fmt.Println(auth)

	
	fmt.Println( reflect.TypeOf( tokenInfo  ) )
	fmt.Println("======= SEBELUM HANDLER USER TOKEN DATA ======= ")
	// fmt.Println( r.Context() )

	return tokenInfo, nil
}

func HandlerUserTokenData (w http.ResponseWriter, r *http.Request) {
	
	fmt.Println("======= HANDLER USER TOKEN DATA ======= ")
	tokenData := r.Context().Value("tokenData").(jwt.MapClaims)
	fmt.Println(tokenData)
	fmt.Println( reflect.TypeOf(tokenData))

	fmt.Println("====== token info ====== ")
	fmt.Println( tokenData["id"])
	fmt.Println( tokenData["username"])
	fmt.Println( tokenData["authority"])

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	type D map[string]interface{}
	Data := D{ "Id": tokenData["id"],  "Username": tokenData["username"], "Authority": tokenData["authority"] }
	respString, _ := json.Marshal( M{ "Message": "Success retrieve User Data", "Status": 200, "userData": Data } )

	w.Write([]byte (respString) )

}


func isLogin(){

}

func isAdmin(){

}

func isSuperUser(){

}

