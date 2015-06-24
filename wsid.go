package main

import (
	"fmt"
	"os"
	"os/user"
	"log"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
)

type thing struct{
 id       sql.NullInt64
 name sql.NullString

}

func main() {
	args := os.Args[1:]

	command := "go"
	if len(args) >= 1 {
		command = args[0]
	}

	db := connect()
	defer db.Close()

	setup(db)
	if command == "add" {
		add(db, args)
	} else if command == "go" {
		gogo(db)
	} else if command == "show" {
		show(db)
	} else if command == "remove" {
		remove(db, args)
	} else {
		fmt.Println("Command not found. Try help.")
	}
}

// adds a thing
func add(db  *sql.DB, args[]string ) {
	if len(args) != 2 {
		fmt.Println("There is nothing to add.")
		return
	}
	name := args[1]

	_, err := db.Exec("insert into things(name) values(?)", name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Added: ", name)
}

// removes a thing
func remove(db  *sql.DB, args[]string) {
	if len(args) != 2 {
		fmt.Println("Not the right amount of argmuents")
		return
	}
	id := args[1]

	_, err := db.Exec("delete from things where id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Removed:", id)
}

func show(db  *sql.DB) {
	rows, err := db.Query("select id, name from things")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		fmt.Println(id, name)
	}
	rows.Close()
}


func gogo(db  *sql.DB) {
	rows, err := db.Query("SELECT * FROM things ORDER BY RANDOM() LIMIT 1;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next(){

		var id int
		var name string
		rows.Scan(&id, &name)
		fmt.Println("You should do:",name)
	} else {
		fmt.Println("There are no things you can do!")
	}

	rows.Close()
}

func  connect() (db  *sql.DB) {
	usr, err := user.Current()
    if err != nil {
        log.Fatal( err )
    }

	db, err = sql.Open("sqlite3", usr.HomeDir + "/wsid.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func setup(db *sql.DB){
	sqlStmt := `
	create table IF NOT EXISTS things (id integer not null primary key, name text);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}