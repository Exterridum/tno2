package protocol

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
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

type Slot struct {
	data map[string]interface{}
}

func (p *Http) actionsPath(ctxPath string, td *model.ThingDescription) []*route {
	var routes []*route
	routes = make([]*route, 0)

	for _, action := range td.Actions {
		slot := Slot{data: make(map[string]interface{})}
		routes = append(routes, p.actionPath(ctxPath, &action, &slot))
		routes = append(routes, p.actionTaskPath(ctxPath, &action, &slot))
	}

	return routes
}

func (p *Http) actionPath(ctxPath string, action *model.Action, slot *Slot) *route {
	return &route{
		name:    action.Name,
		method:  "POST",
		pattern: contextPath(ctxPath, action.Name),
		handlerFunc: func(w http.ResponseWriter, r *http.Request) {
			value := WotObject{}
			json.NewDecoder(r.Body).Decode(&value)

			uuid, _ := newUUID()

			_, rc := p.servers[ctxPath].InvokeAction(
				action.Name,
				value,
				func(status int, text string) {
					slot.data[uuid] = str.Concat("Status: ", status, ", Text: ", text)
				})

			if rc == wot.OK {
				ok(w, uuid)
			} else {
				error(w, "Unknown action.")
			}
		},
	}
}

func (p *Http) actionTaskPath(ctxPath string, action *model.Action, slot *Slot) *route {
	return &route{
		name:    str.Concat(action.Name, "Task"),
		method:  "GET",
		pattern: contextPath(ctxPath, str.Concat(action.Name, "/task/{taskid}")),
		handlerFunc: func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			taskid := vars["taskid"]

			data, isOk := slot.data[taskid]

			if isOk {
				ok(w, data)
			} else {
				error(w, "Task id not found.")
			}

		},
	}
}

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, bool) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", false
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), true
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
