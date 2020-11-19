package main

import (
	"log"
	"net/http"
	// "strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	controllers "go_rest_api/controllers"
)





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
	api.HandleFunc("/user/",controllers.GetAllUser).Methods(http.MethodGet)
	api.HandleFunc("/user/{id}/", controllers.GetById).Methods(http.MethodGet)
	api.HandleFunc("/user/", controllers.CreateUser).Methods(http.MethodPost)
	api.HandleFunc("/user/{id}/", controllers.UpdateUser).Methods(http.MethodPut)
	api.HandleFunc("/user/{id}/", controllers.RemoveUser).Methods(http.MethodDelete)

	// api.HandleFunc("/user", patch).Methods(http.MethodPatch)

	api.HandleFunc("/login/", controllers.LoginUser).Methods(http.MethodPost)
	api.HandleFunc("/verifyToken/", controllers.MiddlewareJWTAuthorization( controllers.HandlerUserTokenData ))


	log.Fatal(http.ListenAndServe(":9090", r))
}



