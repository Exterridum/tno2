package main

import (
	"bytes"
	"container/list"
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

	models := loadModels()

	r := mux.NewRouter().StrictSlash(true)
	appendDevices(r, models)
	// r.PathPrefix("/model").HandlerFunc(Model)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func loadModels() *list.List {
	files, err := ioutil.ReadDir("models")
	if err != nil {
		log.Fatal(err)
	}

	thingModels := list.New()
	for _, fileInfo := range files {
		if fileInfo.IsDir() == false {
			thingModels.PushBack(loadModel(fileInfo.Name()))
		}
	}

	return thingModels
}

const modelRoot = "models"

func loadModel(fileName string) device.Model {
	fullFileName := Concat(modelRoot, "/", fileName)
	file, e := ioutil.ReadFile(fullFileName)

	var model device.Model

	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	} else {
		json.Unmarshal(file, &model)
	}

	return model
}

func appendDevices(r *mux.Router, devices *list.List) {
	deviceUrls := make([]string, devices.Len())

	i := 0
	for e := devices.Front(); e != nil; e = e.Next() {
		model := e.Value.(device.Model)

		r.HandleFunc(Concat("/", model.ID),
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "Device information for -> %s", model.ID)
			})

		deviceUrls[i] = Concat("/", model.ID)
		i++
	}

	var things device.Things
	things.Things = deviceUrls

	r.HandleFunc(("/"),
		func(w http.ResponseWriter, r *http.Request) {
			EncodeJson(w, things)
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
