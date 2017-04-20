package main

import "fmt"
import "net/http"
import "io/ioutil"
import "encoding/json"

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

func (resp response) printinfo() {
	fmt.Printf("Title: %s\nDescription: %s\nAuthor: %s\nAuth: %s\nObjects:\n", resp.Title, resp.Desc, resp.Author, resp.Auth)
	for _, object := range resp.Objects {
		fmt.Printf("\tName: %s\n\tIschecked: %t\n\tDescription: %s\n\n", object.Name, object.Check, object.Desc)
	}
}

func fromjson(src string, v interface{}) error {
	return json.Unmarshal([]byte(src), v)
}

func tojson(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func apihandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	fmt.Println("RAW_RESPONSE", string(body), "\n")
	if err != nil {
		fmt.Println(err)
		return
	}
	var resp response
	fromjson(string(body), &resp)
	fmt.Fprintf(w, "200")
	resp.printinfo()
	fmt.Println()
}

func main() {
	fmt.Println("TODO SERVER\nCopyright by Paul Kramme 2017")
	http.HandleFunc("/api", apihandler)
	http.ListenAndServe(":80", nil)
}
