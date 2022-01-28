package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type JSONFileContent struct {
	Paths []struct{
		Path string 				`json:"path"`
		Get *ReqStructure 	`json:"get"`
		Post *ReqStructure 	`json:"post"`
	} `json:"paths"`
}

type Params []string

type ReqStructure struct {
	Params Params `json:"params"`
	Response json.RawMessage `json:"response"`
}

var r *mux.Router

func main() {
	filePath := flag.String("f", "", "path to file with data to mock")
	address := flag.String("a", ":9090", "address to listen on")
	flag.Parse()

	if flag.NFlag() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.OpenFile(*filePath, os.O_RDONLY, 0444)
	if err != nil {
		log.Fatalf("could not open file '%s': %v", *filePath, err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("could not read file '%s': %v", *filePath, err)
	}

	fileContent := &JSONFileContent{}
	err = json.Unmarshal(b, fileContent)
	if err != nil {
		log.Fatalf("could not unmarshal json file: %v", err)
	}

	r = mux.NewRouter()

	fmt.Println("identified paths:")
	for _, path := range fileContent.Paths {
		if path.Get != nil {
			RegisterNewGetRoute(path.Path, path.Get)
		}
		if path.Post != nil {
			RegisterNewPostRoute(path.Path, path.Post)
		}
	}
	
	http.ListenAndServe(*address, r)
}

func RegisterNewGetRoute(path string, req *ReqStructure) {
	fmt.Println("GET ", path, "?", req.Params)
	query := []string{}
	for _, param := range req.Params {
		query = append(query, param, fmt.Sprintf("{%s}", param))
	}
	r.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		w.Header().Add("Content-Type", "application/json")
		writeResponse(req.Response, w)
	}).Methods(http.MethodGet).Queries(query...)
}

func RegisterNewPostRoute(path string, req *ReqStructure) {
	fmt.Println("POST", path, "?", req.Params)
	query := []string{}
	for _, param := range req.Params {
		query = append(query, param, fmt.Sprintf("{%s}", param))
	}
	r.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		writeResponse(req.Response, w)
	}).Methods(http.MethodPost).Queries(query...)
}

func writeResponse(res json.RawMessage, w io.Writer) {
	j, err := json.Marshal(&res)
	if err != nil {
			panic(err)
	}
	fmt.Fprint(w, string(j))
}