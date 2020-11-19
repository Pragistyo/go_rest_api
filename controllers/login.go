package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"

	// "reflect"
	"io/ioutil"
	"time"

	"encoding/json"
	db "go_rest_api/db"
	helper "go_rest_api/helper" // for hashing bcrypt
	models "go_rest_api/models"

	jwt "github.com/dgrijalva/jwt-go"
)



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
	var u models.User
	row := conn.QueryRow( context.Background(), 
							"SELECT id, username, password, authority FROM Users WHERE username=$1",   
							requestBody.Username )

	err = row.Scan(&u.Id, &u.Username, &u.Password, &u.Authority )
	if err!=nil {
		log.Println("Error in fetch user login: ", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode( models.GetLoginResponse(   " user not found ", 404, "-",0 ) )
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
		json.NewEncoder(w).Encode(models.GetLoginResponse( " password not match ", 401, "-", 0 ) )
		return
	}

	var APPLICATION_NAME = "go_rest_api_App_v1.0"
	// var LOGIN_EXPIRATION_DURATION = time.Duration(1) * time.Hour
	var LOGIN_EXPIRATION_DURATION = time.Minute * 30
	var JWT_SIGNING_METHOD = jwt.SigningMethodHS256
	var JWT_SIGNATURE_KEY = []byte( os.Getenv( "SECRET_KEY_JWT" ) )
	log.Println( reflect.TypeOf(JWT_SIGNATURE_KEY))

	type MyClaims struct {
		jwt.StandardClaims
		Id 			int32		`json:"id"`
		Username	string		`json:"username"`
		Authority	int32		`json:"authority"`
	}

	claims := MyClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer: APPLICATION_NAME,
			ExpiresAt: time.Now().Add(LOGIN_EXPIRATION_DURATION).Unix(),
		} ,
		Id: u.Id,
		Username: u.Username,
		Authority: u.Authority,
	}
	sign := jwt.NewWithClaims(
		JWT_SIGNING_METHOD,
		claims,
	)

	token, err := sign.SignedString(JWT_SIGNATURE_KEY)

	if err!= nil {
		log.Println("Error make credential: ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	
	//update login
	var lastLogin time.Time
	lastLogin = time.Now()

	var sqlStatement string =`
	UPDATE Users 
	SET last_login =$1
	WHERE id = $2
	`
	updLogin, err := conn.Exec(context.Background(), sqlStatement, lastLogin, u.Id)
	if err!= nil {
		log.Panic(err)
	}


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader( http.StatusOK )

	body := models.GetLoginResponse (
		"Success Login!",
		200,
		"Bearer: "+token,
		updLogin.RowsAffected(),
	)

	serializedBody, _ := json.Marshal(body)
	_, _ = w.Write(serializedBody)
	return

	// json.NewEncoder(w).Encode(resp)
}

