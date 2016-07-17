package protocol

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/conas/tno2/util/str"
	"github.com/conas/tno2/wot"
	"github.com/conas/tno2/wot/model"
	"github.com/gorilla/mux"
)

//Based on http://thenewstack.io/make-a-restful-json-api-go/

type Http struct {
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

func NewHttp(port int) *Http {
	// r.PathPrefix("/model").HandlerFunc(Model)
	return &Http{
		port:    port,
		router:  mux.NewRouter().StrictSlash(true),
		hrefs:   make([]string, 0),
		servers: make(map[string]*wot.Server),
	}
}

func (p *Http) Bind(ctxPath string, s *wot.Server) {
	td := s.GetDescription()
	p.servers[ctxPath] = s
	p.createRoutes(ctxPath, td)
	//Update TD uris by created protocol bind
	td.Uris = append(td.Uris, str.Concat("http://localhost:8080", ctxPath))
}

func (p *Http) Start() {
	port := str.Concat(":", strconv.Itoa(p.port))
	log.Fatal(http.ListenAndServe(port, p.router))
}

// ----- ThingDescription parser methods

func (p *Http) createRoutes(ctxPath string, td *model.ThingDescription) {
	var routes []*route
	routes = make([]*route, 0)

	routes = append(routes, p.rootPath(ctxPath, td))
	routes = append(routes, p.descriptionPath(ctxPath, td))

	for _, path := range p.propertiesPath(ctxPath, td) {
		routes = append(routes, path)
	}

	for _, route := range routes {
		p.addRoute(route)
	}
}

func (p *Http) addRoute(route *route) {
	p.router.
		Methods(route.method).
		Path(route.pattern).
		Name(route.name).
		Handler(route.handlerFunc)
}

func (p *Http) rootPath(ctxPath string, td *model.ThingDescription) *route {
	return &route{
		name:    "Index",
		method:  "GET",
		pattern: contextPath(ctxPath, ""),
		handlerFunc: func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Device information for -> %s", td.Name)
		},
	}
}

func (p *Http) descriptionPath(ctxPath string, td *model.ThingDescription) *route {
	return &route{
		name:    "model",
		method:  "GET",
		pattern: contextPath(ctxPath, "description"),
		handlerFunc: func(w http.ResponseWriter, r *http.Request) {
			ok(w, td)
		},
	}
}

func (p *Http) propertiesPath(ctxPath string, td *model.ThingDescription) []*route {
	var routes []*route
	routes = make([]*route, 0)

	for _, prop := range td.Properties {
		routes = append(routes, p.getPropertyPath(ctxPath, &prop))

		if prop.Writable {
			routes = append(routes, p.setPropertyPath(ctxPath, &prop))
		}
	}

	return routes
}

func (p *Http) getPropertyPath(ctxPath string, prop *model.Property) *route {
	e := Encoder(prop)

	return &route{
		name:    prop.Name,
		method:  "GET",
		pattern: contextPath(ctxPath, prop.Name),
		handlerFunc: func(w http.ResponseWriter, r *http.Request) {
			name := prop.Name
			promise, rc := p.servers[ctxPath].GetProperty(name)

			if rc == wot.OK {
				value := promise.Wait()
				ok(w, e(value))
			} else {
				error(w, "Unknown property.")
			}
		},
	}
}

func (p *Http) setPropertyPath(ctxPath string, prop *model.Property) *route {
	d := Decoder(prop)

	return &route{
		name:    prop.Name,
		method:  "PUT",
		pattern: contextPath(ctxPath, prop.Name),
		handlerFunc: func(w http.ResponseWriter, r *http.Request) {
			name := prop.Name
			value := d(r.Body)
			promise, rc := p.servers[ctxPath].SetProperty(name, value)

			if rc == wot.OK {
				promise.Wait()
			} else {
				error(w, "Unknown property.")
			}
		},
	}
}

func contextPath(ctxPath, element string) string {
	return str.Concat(ctxPath, "/", element)
}

func ok(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}

func error(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "text/plain")
	json.NewEncoder(w).Encode(payload)
}
