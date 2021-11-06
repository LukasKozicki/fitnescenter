package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	flag.Parse()

	if _, err := os.Stat("./db/fitcenter.db"); errors.Is(err, os.ErrNotExist) {
		// DB file not exist
		fmt.Print("Database file not found - ")
		if *devEnv == "dev" {
			fmt.Println("Don't be stressed. You are in dev mode")
			fmt.Println("In the next steps, the database file will be created")
		} else {
			fmt.Println("Correct the path or start developer mode")
			fmt.Println("To change the application mode to the run path, add the -v dev parameter")
		}
	} else {
		// DB file exist
		if *devEnv == "dev" {
			// Removing the database
			fmt.Println("The development version has been launched - data will be deleted")
			os.Remove("./db/fitcenter.db")
		} else {
			fmt.Println("The development version is disabled - data will not be deleted")
		}
	}

	db, err := sql.Open("sqlite3", "./db/fitcenter.db")
	Check(err)
	defer db.Close()

	// Import of the database structure and default data - dev mode
	if *devEnv == "dev" {
		ExecuteQuery(db, queryCreateAllTables)
		ExecuteQuery(db, queryClearAllTables)
		ExecuteQuery(db, userTableString)
		ExecuteQuery(db, measurementsUnitsTableString)
		ExecuteQuery(db, measurementsTypesTableString)
	}

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images/"))))
	http.HandleFunc("/yourfit", viewHandler)
	http.HandleFunc("/yourfit/new", newHandler)
	http.HandleFunc("/yourfit/create", createHandler)
	err = http.ListenAndServe("localhost:8081", nil)
	log.Fatal(err)

}
