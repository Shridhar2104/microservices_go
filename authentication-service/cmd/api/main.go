package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/jackc/pgx/v4"
)

const webPort = "80"
var counts int64=0

type Config struct{
	DB *sql.DB
	Models data.Models
}
func main(){
	log.Println("Starting auth service")
	conn := connectToDB()

	if conn==nil{
		log.Panic("Can't connect to database")
	}

	app:= &Config{
		DB:conn,
		Models:data.New(conn),
	}
	srv:= &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),

	}

	err:= srv.ListenAndServe()

	if err!= nil{
		log.Panic(err)
	}
}
func openDB(dsn string) (*sql.DB, error){

	db, err:= sql.Open("pgx", dsn)
	if err!= nil{
		return nil, err
	}

	err = db.Ping()

	if err!= nil{
		return nil, err

	}

	return db, nil
}

func connectToDB() *sql.DB{
	dsn := os.Getenv("DSN")
	for{
		connection , err:= openDB(dsn)
		if err!= nil{
			log.Println("Postgress db not ready yet")
			counts++

		}else{
			log.Println("Connection successful")
			return connection
		}

		if counts>10{
			log.Println(err)
			return nil
		}

		log.Println("Backing off for 2 seconds")

		time.Sleep(2*time.Second)
		continue
	}	
	
}

