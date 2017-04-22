package main

import "fmt"
import "net/http"
import "io/ioutil"
import "encoding/json"
import "flag"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

var info bool

type config struct {
	Api_site_prefix string
	Listen           string
	Info_printing   bool
	Sql_server      string
	Sql_user        string
	Sql_password    string
	Sql_port        int
}

type Object struct {
	Check bool
	Desc  string
	Name  string
}

type response struct {
	Title   string
	Desc    string
	Author  string
	Auth    string
	Objects []Object
}

type User struct {
	Id string
	Pw string
}

func (resp response) printinfo() {
	fmt.Printf("Title: %s\nDescription: %s\nAuthor: %s\nAuth: %s\nObjects:\n", resp.Title, resp.Desc, resp.Author, resp.Auth)
	for _, object := range resp.Objects {
		fmt.Printf("\tName: %s\n\tIschecked: %t\n\tDescription: %s\n", object.Name, object.Check, object.Desc)
	}
	fmt.Println()
}

func fromjson(src string, v interface{}) error {
	return json.Unmarshal([]byte(src), v)
}

func tojson(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func userhandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
}

func main() {
	fmt.Println("TODO SERVER\nCopyright by Paul Kramme 2017")

	var iterator int

	configfile, err := ioutil.ReadFile("./config.json")
	var conf config
	if err != nil {
		fmt.Println("No config.json found in current directory, expecting arguments.")
	} else {
		err = nil
		err = fromjson(string(configfile), &conf)
		if err != nil {
			panic(err)
		}
	}

	infoprinting := flag.Bool("info", false, "Printing incoming api usage and number of connections")
	flag.Parse()

	fmt.Print("Connecting to database: ")
	db_string := fmt.Sprintf("%s:%s@tcp(%s:%d)/todo", conf.Sql_user, conf.Sql_password, conf.Sql_server, conf.Sql_port)
	db, err := sql.Open("mysql", db_string)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Print("success\n")
	fmt.Print("Testing database connection: ")
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Print("success\n")

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "400 %s", err)
			return
		}
		var resp response
		err = fromjson(string(body), &resp)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Get user
		// Compare user with post author
		// post post to database

		fmt.Fprintf(w, "200")

		if *infoprinting == true {
			iterator++
			fmt.Println(iterator)
			resp.printinfo()
		}
	})
	fmt.Println("Initialization complete.")
	http.ListenAndServe(conf.Listen, nil)
}
