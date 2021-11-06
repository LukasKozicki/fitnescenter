package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func viewHandler(writer http.ResponseWriter, request *http.Request) {
	Types := getStringsDB("SELECT measurements_types.name from measurements LEFT JOIN measurements_types on measurements_types.id = measurements.mtype order by measurements.id DESC limit 10")
	Ids := getStringsDB("SELECT id from measurements order by id DESC limit 10")
	dates := getStringsDB("SELECT date from measurements order by id DESC limit 10")
	values := getStringsDB("SELECT value from measurements order by id DESC limit 10")
	html, err := template.ParseFiles("templates/view.html")
	Check(err)
	yourfit := DbData{
		ReadingCount: len(dates),
		ReadingId:    Ids,
		ReadingData:  dates,
		ReadingValue: values,
		ReadingType:  Types,
	}
	err = html.Execute(writer, yourfit)
	Check(err)
}

func newHandler(writer http.ResponseWriter, request *http.Request) {
	NextID, err := strconv.Atoi(strings.Join(getStringsDB("SELECT id from measurements order by id DESC limit 1"), ""))
	Check(err)
	Types := getStringsDB("SELECT name FROM measurements_types order by id DESC limit 10")
	html, err := template.ParseFiles("templates/new.html")
	Check(err)
	yourfit := NewDbData{
		ReadingCount:    (NextID + 1),
		DictionaryTypes: Types,
	}
	err = html.Execute(writer, yourfit)
	Check(err)
}

func createHandler(writer http.ResponseWriter, request *http.Request) {
	NextID := request.FormValue("nextid")
	ReadingValue := request.FormValue("value")
	TypeValue := strings.Join(getStringsDB("SELECT id FROM measurements_types where name = '"+request.FormValue("types")+"' order by id DESC limit 1"), "")
	dt := time.Now()
	time := dt.Format("01-02-2006 15:04:05")
	db, err := sql.Open("sqlite3", "./db/fitcenter.db")
	Check(err)
	defer db.Close()
	var measurementsTableString = "insert into measurements(id, date, value, mtype, user) values(" + NextID + ", '" + time + "', " + ReadingValue + ", " + TypeValue + ", 1)"
	ExecuteQuery(db, measurementsTableString)
	http.Redirect(writer, request, "/yourfit", http.StatusFound)
}

func getStringsDB(query string) []string {
	db, err := sql.Open("sqlite3", "./db/fitcenter.db")
	Check(err)
	defer db.Close()
	var lines []string
	rows, err := db.Query(query)
	Check(err)
	defer rows.Close()
	for rows.Next() {
		var date string
		err = rows.Scan(&date)
		Check(err)
		lines = append(lines, date)
	}
	err = rows.Err()
	Check(err)
	return lines
}

func ExecuteQuery(db *sql.DB, query string) {
	_, err := db.Exec(query)
	Check(err)
}

// func PrepareDb(db *sql.DB, query string) {
// 	_, err := db.Exec(query)
// 	Check(err)
// }
