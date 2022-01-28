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
	"path/filepath"

	"github.com/gorilla/mux"
)

type JSONFileContent struct {
	Paths []struct {
		Path string        `json:"path"`
		Get  *ReqStructure `json:"get"`
		Post *ReqStructure `json:"post"`
	} `json:"paths"`
}

type ReqStructure struct {
	Params           Params          `json:"params"`
	Response         json.RawMessage `json:"response"`
	ResponseFromFile *string         `json:"responseFromFile"`
}

type Params []string

var (
	r             *mux.Router
	fileDirectory string // for holding the directory of the json file so that
	// any relative references can be handled
)

func main() {
	filePath := flag.String("f", "", "path to file with data to mock")
	address := flag.String("a", ":9090", "address to listen on")
	flag.Parse()

	if flag.NFlag() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	fileDirectory = filepath.Dir(*filePath)

	fileContent, err := readFromFile(*filePath)
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
			RegisterNewRoute(http.MethodGet, path.Path, path.Get)
		}
		if path.Post != nil {
			RegisterNewRoute(http.MethodPost, path.Path, path.Post)
		}
	}

	http.ListenAndServe(*address, r)
}

func RegisterNewRoute(method string, path string, req *ReqStructure) {
	fmt.Println(method, path, "?", req.Params)
	query := []string{}
	for _, param := range req.Params {
		query = append(query, param, fmt.Sprintf("{%s}", param))
	}

	response := req.Response

	if req.ResponseFromFile != nil {
		fileContent, err := readFromFile(filepath.Join(fileDirectory, *req.ResponseFromFile))
		if err != nil {
			log.Fatalf("could not read file '%s': %v", *req.ResponseFromFile, err)
		}

		response = fileContent
	}

	r.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(fmt.Sprintf("[HIT | %s]", r.Method), r.URL)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "authorization,content-type")
		writeResponse(response, w)
	}).Methods(method, http.MethodOptions).Queries(query...)
}

// writeResponse writes given message in `res` to the writer `w`
func writeResponse(res json.RawMessage, w io.Writer) {
	j, err := json.Marshal(&res)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(j))
}

// readFromFile reads the contents of given `filePath`
func readFromFile(filePath string) ([]byte, error) {
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0444)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return b, nil
}
