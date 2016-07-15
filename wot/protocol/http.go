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

//Based on http://thenewstack.io/make-a-restful-json-api-go/

var store = make(map[string]interface{})

type ProtoHttp struct {
	port    int
	router  *mux.Router
	hrefs   []string
	servers map[string]*wot.Server
}

type route struct {
	name        string
	method      string
	pattern     string
	handlerFunc http.HandlerFunc
}

// ----- Server API methods

func Http(port int) *ProtoHttp {
	// r.PathPrefix("/model").HandlerFunc(Model)
	return &ProtoHttp{
		port:    port,
		router:  mux.NewRouter().StrictSlash(true),
		hrefs:   make([]string, 0),
		servers: make(map[string]*wot.Server),
	}
}

func (p *ProtoHttp) Bind(ctxPath string, s *wot.Server) {
	td := s.GetDescription()
	p.servers[ctxPath] = s
	p.createRoutes(ctxPath, td)
	//Update TD uris by created protocol bind
	td.Uris = append(td.Uris, Concat("http://localhost:8080", ctxPath))
}

func (p *ProtoHttp) Start() {
	port := Concat(":", strconv.Itoa(p.port))
	log.Fatal(http.ListenAndServe(port, p.router))
}

// ----- ThingDescription parser methods

func (p *ProtoHttp) createRoutes(ctxPath string, td *model.ThingDescription) {
	var routes []*route
	routes = make([]*route, 0)

	routes = append(routes, rootPath(td))
	routes = append(routes, descriptionPath(td))

	for _, path := range propertiesPath(td) {
		routes = append(routes, path)
	}

	for _, route := range routes {
		p.addRoute(ctxPath, route)
	}
}

func (p *ProtoHttp) addRoute(ctxPath string, route *route) {
	p.router.
		Methods(route.method).
		Path(Concat(ctxPath, "/", route.pattern)).
		Name(route.name).
		Handler(route.handlerFunc)
}

func rootPath(td *model.ThingDescription) *route {
	return &route{
		name:    "Index",
		method:  "GET",
		pattern: "",
		handlerFunc: func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Device information for -> %s", td.Name)
		},
	}
}

func descriptionPath(td *model.ThingDescription) *route {
	return &route{
		name:    "model",
		method:  "GET",
		pattern: "description",
		handlerFunc: func(w http.ResponseWriter, r *http.Request) {
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
		name:    prop.Hrefs[0],
		method:  "GET",
		pattern: prop.Hrefs[0],
		handlerFunc: func(w http.ResponseWriter, r *http.Request) {
			createResponse(w, e(store[prop.Hrefs[0]]))
		},
	}
}

func setPropertyPath(prop *model.Property) *route {
	d := Decoder(prop)

	return &route{
		name:    prop.Hrefs[0],
		method:  "PUT",
		pattern: prop.Hrefs[0],
		handlerFunc: func(w http.ResponseWriter, r *http.Request) {
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
