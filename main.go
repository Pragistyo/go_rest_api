package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"encoding/json"

	"github.com/joho/godotenv"
)

type User struct {	
	Id  		int32   	`json:"id"` //,omitempty
	Username 	string    	`json:"username"` //,omitempty
	Password	string    	`json:"password"` //,omitempty
	Authority 	int32    	`json:"authority"` //,omitempty
	Created_on	time.Time    `json:"created_on"` //,omitempty
	Last_login	*time.Time		  `json:"last_login"` //,omitempty
}


func get(w http.ResponseWriter,r *http.Request){
	
	w.WriteHeader(http.StatusOK)

	conn, err := pgx.Connect(context.Background(), os.Getenv("ELEPHANT_URL"))
	if err != nil {
		log.Println(err)
	} else{
		log.Println(conn)
	}
	defer conn.Close(context.Background())
	var u User
	var arr_user []User

	rows, err := conn.Query(context.Background(), "SELECT id,username,password FROM users")
	if err != nil {
        log.Fatal(err)
	}
	
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&u.Id, &u.Username, &u.Password); err != nil {
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

	conn, err := pgx.Connect(context.Background(), os.Getenv("ELEPHANT_URL"))
	if err != nil {
		log.Println(err)
	} else{
		log.Println(conn)
	}
	defer conn.Close(context.Background())

	var u User
	row := conn.QueryRow(
        context.Background(),
        "SELECT * FROM Users WHERE id=$1",   id)
	err = row.Scan(&u.Id, &u.Username, & u.Password, &u.Authority, &u.Created_on, &u.Last_login )
	if err!=nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func post(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write( [ ]byte ( `{"message":"method POST being called"}` ))
}


func put(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type", "applicaiton/json")
	w.WriteHeader(http.StatusOK)
	w.Write([ ]byte (`{"message":"method PUT being called"}` ))
}

func patch(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([ ]byte (`{"message":"method PATCH being called blablablabl"}`))
}

func delete(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([ ]byte ( `{"message":"method DELETE being called"}`))
}

func params(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	pathParams := mux.Vars(r)
	id := pathParams["id"]
	fmt.Println(id)
}



func main(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r:= mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/",get).Methods(http.MethodGet)
	api.HandleFunc("/{id}/", getById).Methods("GET")
	api.HandleFunc("/",post).Methods(http.MethodPost)
	api.HandleFunc("/", put).Methods(http.MethodPut)
	api.HandleFunc("/", patch).Methods(http.MethodPatch)
	api.HandleFunc("/", delete).Methods(http.MethodDelete)
	log.Println(http.ListenAndServe(":9090",r))
	log.Fatal(http.ListenAndServe(":9090", r))
}