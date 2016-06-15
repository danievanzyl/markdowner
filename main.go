package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"log"
	"net/http"
)

var addr = flag.String("addr", "127.0.0.1:3002", "address to listen on")
var css = flag.String("css", "./static/style.css", "path to CSS")

func mdHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("mdHandler: requested")
	// log.Printf("%#v", req)
	vars := mux.Vars(req)
	file := vars["file"]
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		http.Error(w, req.URL.Path, http.StatusNotAcceptable)
	}
	log.Printf("%s", dat)
	output := blackfriday.MarkdownCommon(dat)
	log.Printf("%s", output)

	w.Header().Set("Content-Type", "text/html;charset=UTF-8")
	bo := bufio.NewWriter(w)
	fmt.Fprintf(bo, "<link href=\"%s\" rel=\"stylesheet\"> </link>\n", *css)
	fmt.Fprintf(bo, "%s", output)
	bo.Flush()
}

func handlers() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc(`/{file:[a-zA-Z0-9\-]+.md}`, mdHandler).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	return r
}

func main() {
	flag.Parse()
	log.Printf("Starting srv %s", *addr)
	err := http.ListenAndServe(*addr, handlers())
	if err != nil {
		log.Fatal(err)
	}
}
