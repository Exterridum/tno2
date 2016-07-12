package protocol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/conas/tno2/wot/model"
	"github.com/gorilla/mux"
)

var store = make(map[string]interface{})

//http://thenewstack.io/make-a-restful-json-api-go/
type ProtoHttp struct {
	port   int
	router *mux.Router
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
	}
}

func (p *ProtoHttp) Bind(td *model.ThingDescription) {
	routes := parse(td)

	for _, route := range routes {
		p.append(route)
	}
}

func (p *ProtoHttp) Start() {
	port := Concat(":", strconv.Itoa(p.port))
	log.Fatal(http.ListenAndServe(port, p.router))
}

func (p *ProtoHttp) append(route *route) {
	p.router.
		Methods(route.Method).
		Path(route.Pattern).
		Name(route.Name).
		Handler(route.HandlerFunc)
}

func parse(td *model.ThingDescription) []*route {
	var routes []*route
	routes = make([]*route, 0)

	routes = append(routes, rootPath(td))
	routes = append(routes, modelPath(td))

	for _, path := range propertiesPath(td) {
		routes = append(routes, path)
	}

	return routes
}

func rootPath(td *model.ThingDescription) *route {
	return &route{
		"Index",
		"GET",
		Concat("/", td.Name),
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Device information for -> %s", td.Name)
		},
	}
}

func modelPath(td *model.ThingDescription) *route {
	return &route{
		"model",
		"GET",
		Concat("/", td.Name, "/model"),
		func(w http.ResponseWriter, r *http.Request) {
			EncodeJson(w, td)
		},
	}
}

func propertiesPath(td *model.ThingDescription) []*route {
	var routes []*route
	routes = make([]*route, 0)

	for _, prop := range td.Properties {
		routes = append(routes, getProperty(&prop, td.Name))

		if prop.Writable {
			routes = append(routes, setProperty(&prop, td.Name))
		}
	}

	return routes
}

func getProperty(prop *model.Property, ctx string) *route {
	path := Concat("/", ctx, "/", prop.Hrefs[0])
	e := Encoder(prop)

	store[path] = 5

	return &route{
		path,
		"GET",
		path,
		func(w http.ResponseWriter, r *http.Request) {
			EncodeJson(w, e(w, r))
		},
	}
}

func setProperty(prop *model.Property, ctx string) *route {
	path := Concat("/", ctx, "/", prop.Hrefs[0])
	d := Decoder(prop)

	return &route{
		path,
		"PUT",
		path,
		func(w http.ResponseWriter, r *http.Request) {
			store[path] = d(w, r)
		},
	}
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
