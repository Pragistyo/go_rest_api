package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"reflect"
	// "strconv"
	"encoding/json"

	"github.com/gorilla/mux"
	// "github.com/jackc/pgx/v4"
	pool "github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type User struct {	
	Id  		int32   	`json:"id"` //,omitempty
	Username 	string    	`json:"username"` //,omitempty
	Password	string    	`json:"password"` //,omitempty
	Authority 	int32    	`json:"authority"` //,omitempty
	Created_on	time.Time    `json:"created_on"` //,omitempty
	Last_login	*time.Time		  `json:"last_login"` //,omitempty
}

func dbConnect() *pool.Pool {
	conn, err := pool.Connect(context.Background(), os.Getenv( "ELEPHANT_URL" ))
	if err!= nil {
		log.Fatal(err)
	}
	log.Println(reflect.TypeOf(conn))
	return conn
}


func get(w http.ResponseWriter,r *http.Request){
	
	w.WriteHeader(http.StatusOK)
	
	conn := dbConnect();
	defer conn.Close()
	var u User
	var arr_user []User

	rows, err := conn.Query(context.Background(), "SELECT * FROM users ORDER BY id ASC")
	if err != nil {
        log.Fatal(err)
	}
	
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&u.Id, &u.Username, &u.Password,&u.Authority, &u.Created_on, &u.Last_login);
		 err != nil {
			log.Fatal(err.Error())

		} else {
			arr_user = append(arr_user, u)
		}
	}
	
	fmt.Println(u)
	type Response struct {
		Message string   `json:"message"`
		Status int		  `json:"status"`
		Data []User		`json:"data"`
	}
	var resp Response
	resp.Message = "Success Query GET ALL"
	resp.Status = 200
	resp.Data = arr_user
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(resp)

	// w.Write([]byte (`{"message":"method GET being called"}`) )

}

func getById (w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	id := params["id"]

	conn := dbConnect();
	defer conn.Close()

	var u User
	row := conn.QueryRow( context.Background(), "SELECT * FROM Users WHERE id=$1",   id)

	err := row.Scan(&u.Id, &u.Username, & u.Password, &u.Authority, &u.Created_on, &u.Last_login )
	if err!=nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func post(w http.ResponseWriter,r *http.Request){
	conn := dbConnect();
	defer conn.Close()

	type Response struct {
		Message string  `json:"message"`
		Status int32	`json:"status"`
		New_Id int32		`json:"new_id`
	}

	err := r.ParseMultipartForm(4096)
	if err != nil {
		panic(err)
	}

	username := r.FormValue("username")
	passwordRaw := r.FormValue("password")
	authority := r.FormValue("authority")
	created_on := time.Now()
	
	password, err := HashPassword(passwordRaw)
	if err!=nil {
		log.Fatal(err)
	}


	sqlStatement := `
					INSERT INTO Users (username, password, Authority, created_on, last_login)
					VALUES ($1, $2, $3, $4, $5)
					RETURNING id
					`
	var id int32 = 0
	// log.Println(sqlStatement)
	err = conn.QueryRow(context.Background(), sqlStatement, username, password, authority, created_on, nil).Scan(&id)
	if err != nil {
	  panic(err)
	}

	var resp Response
	resp.Message = "success"
	resp.Status = 201
	resp.New_Id = id


	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
	// w.Write( [ ]byte ( `{"message":"method POST being called"}` ))
}


func put(w http.ResponseWriter,r *http.Request){
	param := mux.Vars(r)
	id := param["id"]

	err := r.ParseMultipartForm(4096)
	if err != nil {
		panic(err)
	}

	username := r.FormValue("username")
	passwordRaw := r.FormValue("password")
	authority := r.FormValue("authority")

	password, err := HashPassword(passwordRaw)
	if err!=nil {
		log.Fatal(err)
	}

	log.Println(authority)

	type Response struct {
		Message string  `json:"message"`
		Status int32	`json:"status"`
		RowAffected int64 `json:"row_affected"`
	}

	conn := dbConnect();
	defer conn.Close()

	var sqlStatement string=`
	UPDATE Users 
	SET username =$1, password= $2, authority = $3
	WHERE id = $4
	`
	resUpd, err := conn.Exec(context.Background(), sqlStatement, username, password, authority, id)
	if err!= nil {
		log.Panic(err)
	}
	var resp Response
	resp.Message = "successfully update"
	resp.Status =200
	resp.RowAffected =resUpd.RowsAffected()

	w.Header().Set("Content-Type", "applicaiton/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
	// w.Write([ ]byte (`{"message":"method PUT being called"}` ))
}



func remove(w http.ResponseWriter,r *http.Request){
	param := mux.Vars(r)
	id := param["id"]

	type Response struct {
		Message string  `json:"message"`
		Status int32	`json:"status"`
		RowAffected int64 `json:"row_affected"`
	}
	conn := dbConnect();
	defer conn.Close()

	// var row_affected int64 = 0
	var sqlStatement string = `
	DELETE from Users WHERE id = $1
	`
	resDel, err := conn.Exec(context.Background(), sqlStatement, id)
	if err!= nil {
		log.Fatal(err)
	}
	log.Println(reflect.TypeOf(resDel.RowsAffected()))
	log.Println(resDel.RowsAffected())

	var resp Response
	resp.Message = fmt.Sprintf("success delete record with id %s", id)
	resp.Status = 200
	resp.RowAffected = resDel.RowsAffected() 


	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(resp)
	// w.Write([ ]byte ( `{"message":"method DELETE being called"}`))
}

func params(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	pathParams := mux.Vars(r)
	id := pathParams["id"]
	fmt.Println(id)
}

func patch(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([ ]byte (`{"message":"method PATCH being called blablablabl"}`))
}

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 11)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err==nil // still dunno what is this
}



func main(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// log.Println(dbConnect())

	r:= mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/",get).Methods(http.MethodGet)
	api.HandleFunc("/{id}/", getById).Methods("GET")
	api.HandleFunc("/",post).Methods(http.MethodPost)
	api.HandleFunc("/{id}/", put).Methods(http.MethodPut)
	api.HandleFunc("/{id}/", remove).Methods(http.MethodDelete)
	api.HandleFunc("/", patch).Methods(http.MethodPatch)
	log.Fatal(http.ListenAndServe(":9090", r))
}