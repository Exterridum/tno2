package protocol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/conas/tno2/wot"
	"github.com/conas/tno2/wot/model"
	"github.com/gorilla/mux"
)

var store = make(map[string]interface{})

//http://thenewstack.io/make-a-restful-json-api-go/
type ProtoHttp struct {
	port   int
	router *mux.Router
	hrefs  []string
}

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func Http(port int) *ProtoHttp {
	// r.PathPrefix("/model").HandlerFunc(Model)
	return &ProtoHttp{
		port,
		mux.NewRouter().StrictSlash(true),
		make([]string, 0),
	}
}

func (p *ProtoHttp) Bind(ctxPath string, s *wot.Server) {
	td := s.GetDescription()
	td.Uris = append(td.Uris, Concat("http://localhost:8080", ctxPath))

	routes := createRoutes(&td)

	for _, route := range routes {
		p.append(ctxPath, route)
	}
}

func (p *ProtoHttp) Start() {
	port := Concat(":", strconv.Itoa(p.port))
	log.Fatal(http.ListenAndServe(port, p.router))
}

func createRoutes(td *model.ThingDescription) []*route {
	var routes []*route
	routes = make([]*route, 0)

	routes = append(routes, rootPath(td))
	routes = append(routes, descriptionPath(td))

	for _, path := range propertiesPath(td) {
		routes = append(routes, path)
	}

	return routes
}

func (p *ProtoHttp) append(ctxPath string, route *route) {
	p.router.
		Methods(route.Method).
		Path(Concat(ctxPath, "/", route.Pattern)).
		Name(route.Name).
		Handler(route.HandlerFunc)
}

func rootPath(td *model.ThingDescription) *route {
	return &route{
		"Index",
		"GET",
		"",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Device information for -> %s", td.Name)
		},
	}
}

func descriptionPath(td *model.ThingDescription) *route {
	return &route{
		"model",
		"GET",
		"description",
		func(w http.ResponseWriter, r *http.Request) {
			createResponse(w, td)
		},
	}
}

func propertiesPath(td *model.ThingDescription) []*route {
	var routes []*route
	routes = make([]*route, 0)

	for _, prop := range td.Properties {
		routes = append(routes, getPropertyPath(&prop))

		if prop.Writable {
			routes = append(routes, setPropertyPath(&prop))
		}
	}

	return routes
}

func getPropertyPath(prop *model.Property) *route {
	e := Encoder(prop)

	return &route{
		prop.Hrefs[0],
		"GET",
		prop.Hrefs[0],
		func(w http.ResponseWriter, r *http.Request) {
			createResponse(w, e(store[prop.Hrefs[0]]))
		},
	}
}

func setPropertyPath(prop *model.Property) *route {
	d := Decoder(prop)

	return &route{
		prop.Hrefs[0],
		"PUT",
		prop.Hrefs[0],
		func(w http.ResponseWriter, r *http.Request) {
			store[prop.Hrefs[0]] = d(r.Body)
		},
	}
}

func createResponse(w http.ResponseWriter, payload interface{}) {
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
