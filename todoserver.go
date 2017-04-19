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
	checkbox bool
	desc     string
}

type response struct {
	Title         string
	Desc          string
	Author        string
	listofobjects []object
}

func fromjson(src string, v interface{}) error {
	return json.Unmarshal([]byte(src), v)
}

func apihandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(w, "RECEIVED POST\n%s\n", body)
	fmt.Println(string(body))
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
