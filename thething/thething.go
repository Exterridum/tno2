package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/conwolab/tno2/thething/device"
	"github.com/gorilla/mux"
)

//http://thenewstack.io/make-a-restful-json-api-go/
func main() {

	file, e := ioutil.ReadFile("./level1.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	// m := new(Dispatch)
	// var m interface{}
	var model device.Model
	json.Unmarshal(file, &model)
	fmt.Printf("Results: %s\n", model.ToString())

	// router := mux.NewRouter().StrictSlash(true)
	// router.HandleFunc("/", Index)
	// router.HandleFunc("/todos", TodoIndex)
	// router.HandleFunc("/todos/{todoID}", TodoShow)

	// log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Todo Index!")
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoID := vars["todoID"]
	fmt.Fprintln(w, "Todo show:", todoID)
}
