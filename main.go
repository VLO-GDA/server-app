package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"google.golang.org/api/script/v1"
)

const scriptID = "MCCJuPe51IX-qdx3z8-7Or6-rz7rtI9Fx"
const newsFeedURL = "http://www.vlo.gda.pl/vlo/?q=rss.xml"

func init() {
	file, err := os.OpenFile("./vloapp.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0744)
	if err != nil {
		panic("could not set up logger")
	}
	log.SetOutput(io.MultiWriter(file, os.Stdout))
	log.SetFlags(log.Ldate | log.Ltime)
}

func main() {
	log.Println("Starting up vloapp")

	client := GetClient("https://www.googleapis.com/auth/spreadsheets", "https://www.googleapis.com/auth/script.external_request")
	// Generate a service object.
	srv, err := script.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve script Client %v", err)
	}

	router := httprouter.New()

	ln := &Proxy{
		Service: srv,
		Script:  scriptID,
		Name:    "getLuckyNumber",
	}
	router.GET("/lucky-number", ln.Handle)

	tth := &Proxy{
		Service: srv,
		Script:  scriptID,
		Name:    "getHours",
	}
	router.GET("/timetable/hours", tth.Handle)

	tt := &Proxy{
		Service: srv,
		Script:  scriptID,
		Name:    "getTimetable",
		Params: map[string]Middleware{
			"group": func(group string) (interface{}, error) {
				//Group validation
				return group, nil
			},
		},
	}
	router.GET("/timetable/group/:group", tt.Handle)

	rq := &Proxy{
		Service: srv,
		Script:  scriptID,
		Name:    "getRandomQuote",
	}

	router.GET("/quote", rq.Handle)

	media := &Proxy{
		Service: srv,
		Script:  scriptID,
		Name:    "getMedia",
	}
	router.GET("/media", media.Handle)

	router.GET("/news", newsHandler)

	log.Fatalln(http.ListenAndServe(":5555", router))
}

func newsHandler(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	resp, err := http.DefaultClient.Get(newsFeedURL)
	if err != nil {
		rw.WriteHeader(500)
		return
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Write([]byte(content))
}
