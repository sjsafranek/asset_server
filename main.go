package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var users_file string

func init() {
	flag.IntVar(&PORT, "p", DEFAULT_PORT, "Server port")
	flag.Parse()
}

func main() {
	var err error

	router := mux.NewRouter()

	// http://www.alexedwards.net/blog/a-recap-of-request-handling

	// TODO
	//  - Create util function for this
	// Static Files
	err = os.MkdirAll(ASSETS_DIRECTORY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	router.PathPrefix("/assets/").Handler(
		http.StripPrefix("/assets/", http.FileServer(
			http.Dir(ASSETS_DIRECTORY))))
	//.end

	// File uploader
	router.Handle("/upload", http.HandlerFunc(FileUploadHandler)).Methods("GET", "POST")
	router.Handle("/api/v1/upload", http.HandlerFunc(FileUploadApiV1Handler)).Methods("POST")
	//.end

	router.Use(LoggingMiddleWare, SetHeadersMiddleWare)

	logger.Infof("Magic happens on port %v...", PORT)
	// err = http.ListenAndServe(fmt.Sprintf(":%v", PORT), TrafficMiddleWare(router))
	err = http.ListenAndServe(fmt.Sprintf(":%v", PORT), router)
	if nil != err {
		panic(err)
	}
}
