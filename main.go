package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"personal/go-proxy-service/app/repository/store_impl"
	"personal/go-proxy-service/app/server"
	"personal/go-proxy-service/app/services"

	_ "github.com/lib/pq"
)

func main(){
	dbConnection, err := OpenDBConnection(
		os.Getenv("db_host"),
		os.Getenv("db_port"),
		os.Getenv("db_user"),
		os.Getenv("db_pass"),
		os.Getenv("db_name"),
		"disable",
		"postgres")
	if err != nil {
		log.Fatal(err)
	}

	defer dbConnection.Close()

	store := store_impl.New(dbConnection)

	service := services.NewService(store)

	srv := server.NewServer(service)

	port := os.Getenv("server_port")

	log.Fatal(http.ListenAndServe(":"+port, srv))

}

func OpenDBConnection(host, port, user, password, dbName, sslMode, dbType string) (*sql.DB, error) {
	//Defining connection string for Postgres
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		dbName,
		sslMode,
	)

	//connection to DB
	DB, err := sql.Open(dbType, psqlInfo)
	if err != nil {
		return nil, err
	}

	// Forcing code to actually open up a connection
	err = DB.Ping()
	if err != nil {
		return nil, err
	}
	// Signaling user about successful connection to DB
	fmt.Printf("Connection to DB %s of type %s was established!", dbName, dbType)
	return DB, nil
}
