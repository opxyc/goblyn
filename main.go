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
	"regexp"
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

var (
	getRouteMap map[string]json.RawMessage
	postRouteMap map[string]json.RawMessage
)

func main() {
	filePath := flag.String("f", "", "path to file with data to mock")
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

	getRouteMap = make(map[string]json.RawMessage, len(fileContent.Paths))
	postRouteMap = make(map[string]json.RawMessage, len(fileContent.Paths))
	
	re := regexp.MustCompile("[{}]+")

	fmt.Println("identified paths:")
	for _, path := range fileContent.Paths {
		p := re.ReplaceAllString(path.Path, "")
		if path.Get != nil {
			getRouteMap[p] = path.Get.Response
			fmt.Println("GET ", p)
		}
		if path.Post != nil {
			postRouteMap[p] = path.Post.Response
			fmt.Println("POST", p)
		}
	}

	http.HandleFunc("/", handler)
	
	http.ListenAndServe(":9090", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	requestPath := r.URL
	requestMethod := r.Method

	w.Header().Add("Content-Type", "application/json")

	switch requestMethod {
	case http.MethodGet:
		if value, ok := getRouteMap[requestPath.Path]; ok {
			writeResponse(value, w)
		} else {
			fmt.Println(value, ok)
			write404(w)
		}
		return
	case http.MethodPost:
		if value, ok := postRouteMap[requestPath.Path]; ok {
			writeResponse(value, w)
		} else {
			write404(w)
		}
		return
	default:
		write404(w)
	}
}

func writeResponse(res json.RawMessage, w io.Writer) {
	j, err := json.Marshal(&res)
	if err != nil {
			panic(err)
	}
	fmt.Fprint(w, string(j))
}

func write404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "404")
}