package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"google.golang.org/api/script/v1"
)

// Middleware validates and converts params for usage in google apps script
type Middleware func(string) (interface{}, error)

// Proxy acts as a relay between client and google apps scripts adding authentication layer
type Proxy struct {
	Service *script.Service
	Script  string
	Name    string
	Params  map[string]Middleware
}

// Handle is a http.Handler for accepting requests at specific endpoint
func (p Proxy) Handle(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	request := &script.ExecutionRequest{Function: p.Name}

	//Validate and convert parameters
	for key, mw := range p.Params {
		value, err := mw(ps.ByName(key))
		if err != nil {
			rw.WriteHeader(400) //Bad request
			return
		}
		request.Parameters = append(request.Parameters, value)
	}

	resp, err := p.Service.Scripts.Run(p.Script, request).Do()
	if err != nil {
		rw.WriteHeader(500) //Internal server error
		log.Println("Got error while conctacting google app execution script ", err, request)
		return
	}

	if resp.Error != nil {
		rw.WriteHeader(500)
		err := resp.Error.Details[0].(map[string]interface{})
		log.Println("Script error message: ", err["errorMessage"])
		//Borrowed from https://developers.google.com/apps-script/guides/rest/quickstart/go
		if err["scriptStackTraceElements"] != nil {
			// There may not be a stacktrace if the script didn't start executing.
			log.Printf("Script error stacktrace:\n")
			for _, trace := range err["scriptStackTraceElements"].([]interface{}) {
				t := trace.(map[string]interface{})
				log.Printf("\t%s: %d\n", t["function"], int(t["lineNumber"].(float64)))
			}
		}
		return
	}

	//TODO: Add option to alter this, probably middleware
	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	//rw.Write([]byte(fmt.Sprintf("%+v", resp.Response)))
	if value, ok := resp.Response.(map[string]interface{})["result"].(string); ok {
		rw.Write([]byte(value))
	} else {
		rw.Write([]byte("null"))
	}
}
