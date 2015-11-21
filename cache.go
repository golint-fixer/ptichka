package main

import (
	"encoding/json"
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
		ioutil.WriteFile(path, jsonBlob, 0644)
	}

	err = json.Unmarshal(jsonBlob, &ids)

	return ids, err
}

func saveCache(path string, ids []string) (error) {
	jsonBlob, err := json.MarshalIndent(ids, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(path, jsonBlob, 0644)

	return err
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
