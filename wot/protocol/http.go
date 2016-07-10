package protocol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	gohttp "net/http"

	"github.com/conas/tno2/wot/model"
	"github.com/gorilla/mux"
)

//http://thenewstack.io/make-a-restful-json-api-go/
type Http struct {
	router *mux.Router
}

func New() *Http {
	// r.PathPrefix("/model").HandlerFunc(Model)
	return &Http{mux.NewRouter().StrictSlash(true)}
}

func (http *Http) Attach(model *model.Model) {
	addRootPath(http.router, model)
	addModelPath(http.router, model)
}

func (http Http) Start(port string) {
	log.Fatal(gohttp.ListenAndServe(port, http.router))
}

const modelRoot = "models"

func addRootPath(r *mux.Router, model *model.Model) {
	r.HandleFunc(Concat("/", model.Name),
		func(w gohttp.ResponseWriter, r *gohttp.Request) {
			fmt.Fprintf(w, "Device information for -> %s", model.Name)
		})
}

func addModelPath(r *mux.Router, model *model.Model) {
	r.HandleFunc(Concat("/", model.Name, "/description"),
		func(w gohttp.ResponseWriter, r *gohttp.Request) {
			EncodeJson(w, model)
		})
}

func EncodeJson(w gohttp.ResponseWriter, payload interface{}) {
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

// func loadModels(r *mux.Router) {
// 	files, err := ioutil.ReadDir("models")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var things device.Things
// 	things.Things = make([]string, 0)

// 	for _, fileInfo := range files {
// 		if fileInfo.IsDir() == false {
// 			model := loadModel(fileInfo.Name())

// 			processModel(r, model)

// 			things.Things = append(things.Things, Concat("/", model.ID))
// 		}
// 	}

// 	r.HandleFunc(("/"),
// 		func(w http.ResponseWriter, r *http.Request) {
// 			EncodeJson(w, things)
// 		})
// }
