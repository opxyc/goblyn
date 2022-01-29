package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

type JSONFileContent struct {
	Paths []struct {
		Path   string        `json:"path"`
		Get    *ReqStructure `json:"get"`
		Post   *ReqStructure `json:"post"`
		Patch  *ReqStructure `json:"patch"`
		Put    *ReqStructure `json:"put"`
		Delete *ReqStructure `json:"delete"`
	} `json:"paths"`
}

type ReqStructure struct {
	Params           Params          `json:"params"`
	Response         json.RawMessage `json:"response"`
	ResponseFromFile *string         `json:"responseFromFile"`
}

type Params []string

var (
	r               *mux.Router
	fileDirectory   string
	delayInResponse *uint
)

func main() {
	filePath := flag.String("f", "", "path to file with data to mock")
	address := flag.String("a", ":9090", "address to listen on")
	delayInResponse = flag.Uint("d", 0, "delay to induce before each response in milliseconds")
	flag.Parse()

	if flag.NFlag() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	fileDirectory = filepath.Dir(*filePath)

	fileContent, err := ReadFromFile(*filePath)
	if err != nil {
		log.Fatalf("could not read file '%s': %v", *filePath, err)
	}

	parsedFileContent := &JSONFileContent{}
	err = json.Unmarshal(fileContent, parsedFileContent)
	if err != nil {
		log.Fatalf("could not unmarshal json file: %v", err)
	}

	r = mux.NewRouter()

	fmt.Println("identified paths:")
	for _, path := range parsedFileContent.Paths {
		if path.Get != nil {
			registerNewRoute(http.MethodGet, path.Path, path.Get)
		}
		if path.Post != nil {
			registerNewRoute(http.MethodPost, path.Path, path.Post)
		}
		if path.Patch != nil {
			registerNewRoute(http.MethodPatch, path.Path, path.Patch)
		}
		if path.Put != nil {
			registerNewRoute(http.MethodPut, path.Path, path.Put)
		}
		if path.Delete != nil {
			registerNewRoute(http.MethodDelete, path.Path, path.Delete)
		}
	}

	http.ListenAndServe(*address, r)
}

func registerNewRoute(method string, path string, req *ReqStructure) {
	fmt.Println(method, path, "?", req.Params)
	query := []string{}
	for _, param := range req.Params {
		query = append(query, param, fmt.Sprintf("{%s}", param))
	}

	response := req.Response

	if req.ResponseFromFile != nil {
		fileContent, err := ReadFromFile(filepath.Join(fileDirectory, *req.ResponseFromFile))
		if err != nil {
			log.Fatalf("could not read file '%s': %v", *req.ResponseFromFile, err)
		}

		response = fileContent
	}

	r.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		log.Println(fmt.Sprintf("[HIT | %s]", r.Method), r.URL)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "authorization,content-type")
		time.Sleep(time.Duration(*delayInResponse) * time.Millisecond)
		WriteResponse(response, w)
	}).Methods(method, http.MethodOptions).Queries(query...)
}
