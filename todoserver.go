package main

import "fmt"
import "net/http"
import "io/ioutil"
import "encoding/json"

type config struct {
	Site_prefix string
	Port int
	Sql_server_link string
	
}

type object struct {
	Checkbox bool
	Desc     string
}

type response struct {
	Title   string
	Desc    string
	Author  string
	Auth string
	Objects []object
}

func (resp response) printinfo() {
	fmt.Printf("Title: %s\nDescription: %s\nAuthor: %s\nAuth: %s\n", resp.Title, resp.Desc, resp.Author, resp.Auth)
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
	fmt.Println("RAW_RESPONSE", string(body))
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
