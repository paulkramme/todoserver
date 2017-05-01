package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	"log"
)

var (
	version string = "0.2.1"
)

type config struct {
	Api_site_prefix string
	Listen          string
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

type GithubReleasesAssetsApiResponse struct {
	Browser_download_url string
}

type GithubReleasesApiResponse struct {
	Tag_name string
	Html_url string
	Assets   []GithubReleasesAssetsApiResponse
}

func checkupdate(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("There was a problem with autoupdate:")
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var ghapiresp GithubReleasesApiResponse
	err = fromjson(string(body), &ghapiresp)
	if ghapiresp.Tag_name != version {
		fmt.Printf("Version %s available at %s\n", ghapiresp.Tag_name, ghapiresp.Html_url)
	}
}

func main() {
	fmt.Printf("TODO SERVER %s - Copyright by Paul Kramme 2017\n", version)
	fmt.Println("Licensed under MIT License")

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

	go checkupdate("https://api.github.com/repos/paulkramme/todoserver/releases/latest")

	fmt.Print("Connecting to database: ")
	db_string := fmt.Sprintf("%s:%s@tcp(%s:%d)/todo", conf.Sql_user, conf.Sql_password, conf.Sql_server, conf.Sql_port)
	db, err := sql.Open("mysql", db_string)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("success")

	fmt.Print("Testing database connection: ")
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("success")

	stmt, err := db.Prepare("INSERT INTO todos(name, description, username, objects) VALUES(?,?,?,?)")
	if err != nil {
		log.Println(err)
	}

	http.HandleFunc("/api/add", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			if *infoprinting == true {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "{\"message\": \"%s\"}", err)
			return
		}
		var resp response
		err = fromjson(string(body), &resp)
		if err != nil {
			if *infoprinting == true {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "{\"message\":\"%s\"}", err)

			return
		}

		jsonstmt, _ := tojson(resp.Objects)
		_, err = stmt.Exec(resp.Title, resp.Desc, resp.Author, string(jsonstmt))
		if err != nil {
			log.Println(err)
		}

		// Get user
		// Compare user with post author
		// post post to database

		w.WriteHeader(http.StatusOK)

		if *infoprinting == true {
			iterator++
			fmt.Println(iterator)
			resp.printinfo()
		}
	})

	fmt.Println("Initialization complete.")
	http.ListenAndServe(conf.Listen, nil)
}

func fromjson(src string, v interface{}) error {
	return json.Unmarshal([]byte(src), v)
}

func tojson(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (resp response) printinfo() {
	log.Printf("Title: %s\nDescription: %s\nAuthor: %s\nAuth: %s\nObjects:\n", resp.Title, resp.Desc, resp.Author, resp.Auth)
	for _, object := range resp.Objects {
		fmt.Printf("\tName: %s\n\tIschecked: %t\n\tDescription: %s\n", object.Name, object.Check, object.Desc)
	}
	fmt.Println()
}
