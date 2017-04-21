package main

import "fmt"
import "net/http"
import "io/ioutil"
import "encoding/json"
import "flag"
import "database/sql"
import "github.com/go-sql-driver/mysql"

var iterator int
var info bool

type config struct {
	Site_prefix     string
	Port            int
	Sql_server_link string
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

func apihandler(w http.ResponseWriter, r *http.Request) {
	//defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	//fmt.Println("RAW_RESPONSE", string(body), "\n")
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
	fmt.Fprintf(w, "200")
	
	if info == true {
		iterator++
		fmt.Println(iterator)
		resp.printinfo()
	}
	r.Body.Close()
}

func main() {
	fmt.Println("TODO SERVER\nCopyright by Paul Kramme 2017")

	infoprinting := flag.Bool("info", false, "Printing incoming api usage and number of connections")
	flag.Parse()
	info = *infoprinting

	http.HandleFunc("/api", apihandler)
	http.ListenAndServe(":80", nil)
}
