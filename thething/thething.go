package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/conas/tno2/thething/device"
	"github.com/gorilla/mux"
)

//http://thenewstack.io/make-a-restful-json-api-go/
func main() {
	r := mux.NewRouter().StrictSlash(true)

	loadModels(r)

	// r.PathPrefix("/model").HandlerFunc(Model)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func loadModels(r *mux.Router) {
	files, err := ioutil.ReadDir("models")
	if err != nil {
		log.Fatal(err)
	}

	var things device.Things
	things.Things = make([]string, 0)

	for _, fileInfo := range files {
		if fileInfo.IsDir() == false {
			model := loadModel(fileInfo.Name())

			processModel(r, model)

			things.Things = append(things.Things, Concat("/", model.ID))
		}
	}

	r.HandleFunc(("/"),
		func(w http.ResponseWriter, r *http.Request) {
			EncodeJson(w, things)
		})
}

const modelRoot = "models"

func loadModel(fileName string) device.Model {
	fullFileName := Concat(modelRoot, "/", fileName)
	file, e := ioutil.ReadFile(fullFileName)

	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	var model device.Model

	json.Unmarshal(file, &model)

	return model
}

func processModel(r *mux.Router, model device.Model) {
	addRootPath(r, model)
	addModelPath(r, model)
}

func addRootPath(r *mux.Router, model device.Model) {
	r.HandleFunc(Concat("/", model.ID),
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Device information for -> %s", model.ID)
		})
}

func addModelPath(r *mux.Router, model device.Model) {
	r.HandleFunc(Concat("/", model.ID, "/model"),
		func(w http.ResponseWriter, r *http.Request) {
			EncodeJson(w, model)
		})
}

func EncodeJson(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}

func Concat(strings ...string) string {
	var buffer bytes.Buffer

	for _, str := range strings {
		buffer.WriteString(str)
	}

	return buffer.String()
}
