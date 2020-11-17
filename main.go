package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	// "os"
	"time"
	"reflect"
	// "strconv"
	"encoding/json"
	
	

	"github.com/gorilla/mux"
	// "github.com/jackc/pgx/v4"
	// pool "github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	// "golang.org/x/crypto/bcrypt"
	helper "go_rest_api/helper"  // for hash bcrypt password
	controllers "go_rest_api/controllers"
	db "go_rest_api/db"
	
)

type User struct {	
	Id  		int32   		  `json:"id"` //,omitempty
	Username 	string    		  `json:"username"` //,omitempty
	Password	string    		  `json:"password"` //,omitempty
	Authority 	int32    		  `json:"authority"` //,omitempty
	Created_on	time.Time    	  `json:"created_on"` //,omitempty
	Last_login	*time.Time		  `json:"last_login"` //,omitempty
}



func get(w http.ResponseWriter,r *http.Request){
	
	w.WriteHeader(http.StatusOK)
	
	conn := db.Connect();
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


}

func getById(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	id := params["id"]

	conn := db.Connect();
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
	conn := db.Connect();
	defer conn.Close()

	type Response struct {
		Message string  `json:"message"`
		Status int32	`json:"status"`
		NewId int32	`json:"new_id`
	}

	err := r.ParseMultipartForm(4096)
	if err != nil {
		panic(err)
	}

	username := r.FormValue("username")
	passwordRaw := r.FormValue("password")
	authority := r.FormValue("authority")
	created_on := time.Now()
	
	password, err := helper.HashPassword(passwordRaw)
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
	resp.NewId = id


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

	password, err := helper.HashPassword(passwordRaw)
	if err!=nil {
		log.Fatal(err)
	}

	log.Println(authority)

	type Response struct {
		Message string  `json:"message"`
		Status int32	`json:"status"`
		RowAffected int64 `json:"row_affected"`
	}

	conn := db.Connect();
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

	conn := db.Connect();
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
}


func patch(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([ ]byte (`{"message":"method PATCH being called"}`))
}



func main(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r:= mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/user/",get).Methods(http.MethodGet)
	api.HandleFunc("/user/{id}/", getById).Methods(http.MethodGet)
	api.HandleFunc("/user/",post).Methods(http.MethodPost)
	api.HandleFunc("/user/{id}/", put).Methods(http.MethodPut)
	api.HandleFunc("/user/{id}/", remove).Methods(http.MethodDelete)
	// api.HandleFunc("/user", patch).Methods(http.MethodPatch)
	api.HandleFunc("/login/", controllers.LoginUser).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":9090", r))
}