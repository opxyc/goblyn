package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// writeResponse writes given message in `res` to the writer `w`
func WriteResponse(res json.RawMessage, w io.Writer) {
	j, err := json.Marshal(&res)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(j))
}

// readFromFile reads the contents of given `filePath`
func ReadFromFile(filePath string) ([]byte, error) {
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
