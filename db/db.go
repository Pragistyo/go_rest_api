package server

import (
	"os"
	"reflect"
	"log"
	"context"

	pool "github.com/jackc/pgx/v4/pgxpool"

)

func Connect() *pool.Pool {
	conn, err := pool.Connect(context.Background(), os.Getenv( "ELEPHANT_URL" ))
	if err!= nil {
		log.Fatal(err)
	}
	log.Println(reflect.TypeOf(conn))
	return conn
}