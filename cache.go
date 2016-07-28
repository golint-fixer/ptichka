package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func loadCache(path string) ([]string, error) {
	var ids []string

	pathExists, err := pathExists(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on pathExists(%s): %v", path, err)
		os.Exit(1)
	}

	var jsonBlob []byte
	if pathExists {
		jsonBlob, err = ioutil.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error on ReadFile(%s): %v", path, err)
			os.Exit(1)
		}
	} else {
		jsonBlob = []byte(`[]`)
		if err := ioutil.WriteFile(path, jsonBlob, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error on WriteFile(%s): %v", path, err)
			os.Exit(1)
		}
	}

	err = json.Unmarshal(jsonBlob, &ids)

	return ids, err
}

func saveCache(path string, ids []string) error {
	jsonBlob, err := json.MarshalIndent(ids, "", "  ")
	ifError(err, "Error on MarshalIndent: %s")

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
