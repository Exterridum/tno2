package protocol

import (
	"log"
	"net/http"

	"github.com/conas/tno2/wot/model"
	"github.com/gorilla/mux"
)

//http://thenewstack.io/make-a-restful-json-api-go/
func main() {
	r := mux.NewRouter().StrictSlash(true)

	// loadModels(r)

	// r.PathPrefix("/model").HandlerFunc(Model)

	log.Fatal(http.ListenAndServe(":8080", r))
}

type Http struct {
}

func New() Http {
	// r := mux.NewRouter().StrictSlash(true)
	return Http{}
}

func Start(port int16) {
	// log.Fatal(http.ListenAndServe(":8080", r))
}

func (http Http) Attach(model *device.Model) {

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

// const modelRoot = "models"

// func processModel(r *mux.Router, model device.Model) {
// 	addRootPath(r, model)
// 	addModelPath(r, model)
// }

// func addRootPath(r *mux.Router, model device.Model) {
// 	r.HandleFunc(Concat("/", model.ID),
// 		func(w http.ResponseWriter, r *http.Request) {
// 			fmt.Fprintf(w, "Device information for -> %s", model.ID)
// 		})
// }

// func addModelPath(r *mux.Router, model device.Model) {
// 	r.HandleFunc(Concat("/", model.ID, "/model"),
// 		func(w http.ResponseWriter, r *http.Request) {
// 			EncodeJson(w, model)
// 		})
// }

// func EncodeJson(w http.ResponseWriter, payload interface{}) {
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(payload)
// }

// func Concat(strings ...string) string {
// 	var buffer bytes.Buffer

// 	for _, str := range strings {
// 		buffer.WriteString(str)
// 	}

// 	return buffer.String()
// }
