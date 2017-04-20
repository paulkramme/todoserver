package main

import "fmt"
import "net/http"
import "io/ioutil"
import "encoding/json"

type Site struct {
	Title string
	Body  []byte
}

type object struct {
	Checkbox bool
	Desc     string
}

type response struct {
	Title         string
	Desc          string
	Author        string
	Listofobjects []object
}

func (resp response) printinfo() {
	fmt.Printf("Title: %s\nDescription: %s\nAuthor: %s\n", resp.Author, resp.Desc, resp.Author)
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
	if err != nil {
		fmt.Println(err)
		return
	}
	var resp response
	fromjson(string(body), &resp)
	fmt.Fprintf(w, "RECEIVED POST\n%s\n", resp)
	resp.printinfo()
	fmt.Println()
}

func roothandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Bendix")
}

func main() {
	fmt.Println("TODO SERVER\nCopyright by Paul Kramme 2017")
	http.HandleFunc("/api", apihandler)
	http.HandleFunc("/", roothandler)
	http.ListenAndServe(":80", nil)
}
