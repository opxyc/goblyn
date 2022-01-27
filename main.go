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
		Path string 		`json:"path"`
		Get *Response 		`json:"get"`
		Post *Response 	`json:"post"`
	} `json:"paths"`
}

type Response struct {
	Response json.RawMessage `json:"response"`
}

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

	r := mux.NewRouter()

	fmt.Println("identified paths:")
	for _, path := range fileContent.Paths {
		p := path.Path
		if path.Get != nil {
			fmt.Println("GET ", p)
			r.HandleFunc(path.Path, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "application/json")
				writeResponse(path.Get.Response, w)
			}).Methods(http.MethodGet)
		}
		if path.Post != nil {
			fmt.Println("POST", p)
			r.HandleFunc(path.Path, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "application/json")
				writeResponse(path.Post.Response, w)
			}).Methods(http.MethodPost)
		}
	}
	
	http.ListenAndServe(*address, r)
}

func writeResponse(res json.RawMessage, w io.Writer) {
	j, err := json.Marshal(&res)
	if err != nil {
			panic(err)
	}
	fmt.Fprint(w, string(j))
}