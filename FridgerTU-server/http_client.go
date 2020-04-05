package main

import (
	"io/ioutil"
	"net/http"
)

func executeGetRequest(url string) (res []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return r, nil
}
