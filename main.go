package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// The variable that defines development processes
var devEnv = flag.String("v", "live", "help message for flag n")

var queryCreateAllTables = `
create table users (id integer not null primary key, name text, create_date text);
create table measurements_types (id integer not null primary key, name text, unit integer);
create table measurements_units(id integer not null primary key, name text, symbol text, number integer);
create table measurements (id integer not null primary key, date text, value real, mtype integer not null, user integer not null);
`
var queryClearAllTables = `
delete from users;
delete from measurements_types;
delete from measurements_units;
delete from measurements;
`

var userTableString = "insert into users(id, name) values(1, 'sysadmin')"
var measurementsUnitsTableString = "insert into measurements_units(id, name, symbol, number) values(1, 'mg/dL', 'mg/dL', 0), (2, 'Percent', '%', 2), (3, 'kilogram', 'kg', 1)"
var measurementsTypesTableString = "insert into measurements_types(id, name, unit) values(1, 'glucose', 1), (2, 'weight',3), (3, 'body fat',2)"

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type DbData struct {
	ReadingCount int
	ReadingId    []string
	ReadingData  []string
	ReadingValue []string
	ReadingType  []string
	// ReadingUser  []int
}

type NewDbData struct {
	ReadingCount int
	//ReadingId       []string
	DictionaryTypes []string
}

func getStringsDB(query string) []string {
	db, err := sql.Open("sqlite3", "./db/fitcenter.db")
	check(err)
	defer db.Close()
	var lines []string
	rows, err := db.Query(query)
	check(err)
	defer rows.Close()
	for rows.Next() {
		var date string
		err = rows.Scan(&date)
		check(err)
		lines = append(lines, date)
	}
	err = rows.Err()
	check(err)
	return lines
}

func viewHandler(writer http.ResponseWriter, request *http.Request) {
	//signatures := getStrings("signatures.txt")
	Types := getStringsDB("SELECT measurements_types.name from measurements LEFT JOIN measurements_types on measurements_types.id = measurements.mtype order by measurements.id DESC limit 10")
	Ids := getStringsDB("SELECT id from measurements order by id DESC limit 10")
	dates := getStringsDB("SELECT date from measurements order by id DESC limit 10")
	values := getStringsDB("SELECT value from measurements order by id DESC limit 10")
	html, err := template.ParseFiles("templates/view.html")
	check(err)
	yourfit := DbData{
		ReadingCount: len(dates),
		ReadingId:    Ids,
		ReadingData:  dates,
		ReadingValue: values,
		ReadingType:  Types,
	}
	err = html.Execute(writer, yourfit)
	check(err)
}

func newHandler(writer http.ResponseWriter, request *http.Request) {
	NextID, err := strconv.Atoi(strings.Join(getStringsDB("SELECT id from measurements order by id DESC limit 1"), ""))
	check(err)
	Types := getStringsDB("SELECT name FROM measurements_types order by id DESC limit 10")
	html, err := template.ParseFiles("templates/new.html")
	check(err)
	yourfit := NewDbData{
		ReadingCount:    (NextID + 1),
		DictionaryTypes: Types,
	}
	err = html.Execute(writer, yourfit)
	check(err)
}

func createHandler(writer http.ResponseWriter, request *http.Request) {
	NextID := request.FormValue("nextid")
	ReadingValue := request.FormValue("value")
	TypeValue := strings.Join(getStringsDB("SELECT id FROM measurements_types where name = '"+request.FormValue("types")+"' order by id DESC limit 1"), "")
	dt := time.Now()
	time := dt.Format("01-02-2006 15:04:05")
	db, err := sql.Open("sqlite3", "./db/fitcenter.db")
	check(err)
	defer db.Close()
	var measurementsTableString = "insert into measurements(id, date, value, mtype, user) values(" + NextID + ", '" + time + "', " + ReadingValue + ", " + TypeValue + ", 1)"
	InsertData(db, measurementsTableString)
	http.Redirect(writer, request, "/yourfit", http.StatusFound)
}

func main() {
	flag.Parse()
	dt := time.Now()
	time := dt.Format("01-JAN-2006 15:04:05")
	fmt.Println(time)
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
	check(err)
	defer db.Close()

	// Import of the database structure and default data - dev mode
	if *devEnv == "dev" {
		PrepareDb(db, queryCreateAllTables)
		PrepareDb(db, queryClearAllTables)
		InsertData(db, userTableString)
		InsertData(db, measurementsUnitsTableString)
		InsertData(db, measurementsTypesTableString)
	}
	var id int
	last_id_query := "SELECT id from measurements order by id DESC limit 1"
	rows, err := db.Query(last_id_query)
	check(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id)
		check(err)
	}
	err = rows.Err()
	check(err)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images/"))))
	http.HandleFunc("/yourfit", viewHandler)
	http.HandleFunc("/yourfit/new", newHandler)
	http.HandleFunc("/yourfit/create", createHandler)
	err = http.ListenAndServe("localhost:8081", nil)
	log.Fatal(err)

}

func SelectData(db *sql.DB, query string) {
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var date string
		var value float32
		var mtype int
		var user int
		err = rows.Scan(&id, &date, &value, &mtype, &user)
		check(err)
		fmt.Println(id, date, value, mtype, user)
	}
	err = rows.Err()
	check(err)
}

func InsertData(db *sql.DB, query string) {
	_, err := db.Exec(query)
	check(err)
}

func PrepareDb(db *sql.DB, query string) {
	_, err := db.Exec(query)
	check(err)
}
