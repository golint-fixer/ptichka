package main

import (
	"encoding/json"
	// "github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"log"
	"os"
)

func loadCache(path string) ([]string, error) {
	var ids []string

	pathExists, err := pathExists(path)
	if err != nil {
		log.Fatal(err)
	}

	var jsonBlob []byte
	if pathExists {
		jsonBlob, err = ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		jsonBlob = []byte(`[]`)
		// prettyJsonBlob, err := json.MarshalIndent(jsonBlob, "", "    ")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		ioutil.WriteFile(path, jsonBlob, 0644)
	}

	err = json.Unmarshal(jsonBlob, &ids)

	return ids, err
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
