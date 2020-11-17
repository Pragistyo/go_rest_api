package controllers

import(
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	// "reflect"
	"io/ioutil"
	"time"


	"encoding/json"
	jwt "github.com/dgrijalva/jwt-go"
	helper "go_rest_api/helper"  // for hashing bcrypt
	db "go_rest_api/db"
	models "go_rest_api/models"

)


type User struct {	
	Id  		int32   		  `json:"id"` //,omitempty
	Username 	string    		  `json:"username"` //,omitempty
	Password	string    		  `json:"password"` //,omitempty
	Authority 	int32    		  `json:"authority"` //,omitempty
	Created_on	time.Time    	  `json:"created_on"` //,omitempty
	Last_login	*time.Time		  `json:"last_login"` //,omitempty
}



func LoginUser(w http.ResponseWriter,r *http.Request){
	

	rawReqBody, err := ioutil.ReadAll(r.Body)

	if err!= nil {
		log.Println("====== Print error ioutil read r.Body: ", err)
		w.WriteHeader( http.StatusBadRequest )
		return
	}
	
	type ReqBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	
	requestBody := &ReqBody{}
	
	err = json.Unmarshal(rawReqBody, requestBody)
	if err != nil {
		log.Panicln(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("_+_+_+_+_+_+_+_+_+_+_+_")
	log.Println(requestBody)
	
	
	conn := db.Connect();
	defer conn.Close()

	//checking user exist
	var u User
	row := conn.QueryRow( context.Background(), 
							"SELECT id, username, password, authority FROM Users WHERE username=$1",   
							requestBody.Username )

	err = row.Scan(&u.Id, &u.Username, &u.Password, &u.Authority )
	if err!=nil {
		log.Println("Error in fetch user login: ", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode( models.GetLoginResponse(   " user not found ", 404, "-" ) )
		return
	}

	// cek password

	var password string = requestBody.Password
	var passwordHash string = u.Password
	log.Println ("=== Hash Password: ", passwordHash)
	log.Println ("=== password req.body: ", password)
	
	matchPass := helper.CheckPasswordHash( password, passwordHash)

	fmt.Println(" ======== PASSWORD CHECKING MATCH ======== ")
	fmt.Println("Match: ", matchPass)
	if !matchPass {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.GetLoginResponse( " password not match ", 401, "-" ) )
		return
	}
	

	sign := jwt.New(jwt.GetSigningMethod("HS256"))
	claims := sign.Claims.(jwt.MapClaims)
	claims["id"]= u.Id
	claims["username"]=u.Username
	claims["authority"]=u.Authority
	token, err := sign.SignedString([]byte( os.Getenv( "SECRET_KEY_JWT" ) ))

	if err!= nil {
		log.Println("Error make credential: ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader( http.StatusOK )

	

	body := models.GetLoginResponse (
		"Hello world from chi!",
		200,
		"Authentication Bearer: "+token,
		// Data:   requestBody,
	)

	serializedBody, _ := json.Marshal(body)
	_, _ = w.Write(serializedBody)
	return

	// json.NewEncoder(w).Encode(resp)
}

func isLogin(){

}

func isAdmin(){

}

func isSuperUser(){

}
