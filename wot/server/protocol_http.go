package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/conas/tno2/util/sec"
	"github.com/conas/tno2/util/str"
	"github.com/conas/tno2/wot/model"
	"github.com/gorilla/mux"
)

//Based on http://thenewstack.io/make-a-restful-json-api-go/

type Http struct {
	port       int
	router     *mux.Router
	hrefs      []string
	wotServers map[string]*WotServer
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
		port:       port,
		router:     mux.NewRouter().StrictSlash(true),
		hrefs:      make([]string, 0),
		wotServers: make(map[string]*WotServer),
	}
}

func (p *Http) Bind(ctxPath string, s *WotServer) {
	td := s.GetDescription()
	p.wotServers[ctxPath] = s
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

	for _, path := range p.actionsPath(ctxPath, td) {
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
			sendOK(w, td)
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
			promise, rc := p.wotServers[ctxPath].GetProperty(name)

			if rc == OK {
				value := promise.Wait()
				sendOK(w, e(value))
			} else {
				sendERR(w, rc)
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
			promise, rc := p.wotServers[ctxPath].SetProperty(name, value)

			if rc == OK {
				promise.Wait()
			} else {
				sendERR(w, rc)
			}
		},
	}
}

func (p *Http) actionsPath(ctxPath string, td *model.ThingDescription) []*route {
	var routes []*route
	routes = make([]*route, 0)

	for _, action := range td.Actions {
		actionTaskPaths := make(map[string]*atomic.Value)
		routes = append(routes, p.actionPath(ctxPath, &action, actionTaskPaths))
		routes = append(routes, p.actionTaskPath(ctxPath, &action, actionTaskPaths))
	}

	return routes
}

func (p *Http) actionPath(ctxPath string, action *model.Action, actionTaskPaths map[string]*atomic.Value) *route {
	return &route{
		name:    action.Name,
		method:  "POST",
		pattern: contextPath(ctxPath, action.Name),
		handlerFunc: func(w http.ResponseWriter, r *http.Request) {
			value := WotObject{}
			json.NewDecoder(r.Body).Decode(&value)

			uuid, _ := sec.UUID4()
			ash := newActionStatusHandler()
			actionTaskPaths[uuid] = ash.Value

			_, rc := p.wotServers[ctxPath].InvokeAction(action.Name, value, &ash)

			if rc == OK {
				sendOK(w, url(r, uuid))
			} else {
				sendERR(w, rc)
			}
		},
	}
}

func url(r *http.Request, relativeLink string) *Links {
	linkString := str.Concat("http://", r.Host, r.URL, "/", relativeLink)

	link := Link{
		Rel:  "taskid",
		Href: linkString,
	}

	return &Links{
		Links: append(make([]Link, 0), link),
	}
}

func (p *Http) actionTaskPath(ctxPath string, action *model.Action, actionTaskPaths map[string]*atomic.Value) *route {
	return &route{
		name:    str.Concat(action.Name, "Task"),
		method:  "GET",
		pattern: contextPath(ctxPath, str.Concat(action.Name, "/{taskid}")),
		handlerFunc: func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			taskid := vars["taskid"]

			value, rc := actionTaskPaths[taskid]

			if rc {
				sendOK(w, value.Load())
			} else {
				sendERR(w, rc)
			}
		},
	}
}

func contextPath(ctxPath, element string) string {
	return str.Concat(ctxPath, "/", element)
}

func sendOK(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}

func sendERR(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "text/plain")
	json.NewEncoder(w).Encode(payload)
}

type Links struct {
	Links []Link `json:"links"`
}

type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}
