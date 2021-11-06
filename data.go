package main

import (
	"flag"
	"strconv"
)

// The variable that defines development processes
var devEnv = flag.String("v", "live", "help message for flag n")

var datalimit = strconv.Itoa(10)
var dbLocation = "./db/fitcenter.db"
var dbType = "sqlite3"

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
