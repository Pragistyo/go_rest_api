package controllers

import(
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"reflect"
	"encoding/json"

	db "go_rest_api/db"
	models "go_rest_api/models"
	helper "go_rest_api/helper"  // for hash bcrypt password
	"github.com/gorilla/mux"
	
)

func GetAllUser(w http.ResponseWriter,r *http.Request){
	
	w.WriteHeader(http.StatusOK)
	
	conn := db.Connect()
	defer conn.Close()

	var u models.User
	var arr_user []models.User
	
	log.Println( reflect.TypeOf( u ) )
	log.Println( reflect.TypeOf( arr_user ) )

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
		Message string   			`json:"message"`
		Status int		            `json:"status"`
		Data []models.User  `json:"data"`
	}

	var resp Response
	resp.Message = "Success Query GET ALL"
	resp.Status = 200
	resp.Data = arr_user

	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(resp)


}

func GetById(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	id := params["id"]

	conn := db.Connect();
	defer conn.Close()

	var u  models.User
	row := conn.QueryRow( context.Background(), "SELECT * FROM Users WHERE id=$1",   id)

	err := row.Scan(&u.Id, &u.Username, & u.Password, &u.Authority, &u.Created_on, &u.Last_login )
	if err!=nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func CreateUser(w http.ResponseWriter,r *http.Request){
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


func UpdateUser(w http.ResponseWriter,r *http.Request){
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
	resp.RowAffected = resUpd.RowsAffected()

	w.Header().Set("Content-Type", "applicaiton/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
	// w.Write([ ]byte (`{"message":"method PUT being called"}` ))
}



func RemoveUser(w http.ResponseWriter,r *http.Request){
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