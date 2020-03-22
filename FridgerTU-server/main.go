package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const baseUrl = "http://www.recipepuppy.com/api/?"

func handler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	_, err := writer.Write([]byte(`{"testKey":"` + request.URL.String() + `"}`))
	if err != nil {
		fmt.Println(err)
	}
}

func executeRequest(url string) (response *Response, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := Response{}
	err = json.Unmarshal(r, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func main() {
	resp, err := executeRequest("http://www.recipepuppy.com/api/?i=onions,garlic&q=omelet&p=3")
	if err != nil {
		fmt.Println(err)
		return
	}
	var rr [][]string
	for _, recipe := range resp.Recipes {
		rr = append(rr, strings.Split(recipe.Ingredients, ", "))
	}
	fmt.Printf("%v\n", rr)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	err = http.ListenAndServe(":8000", mux)
	if err != nil {
		fmt.Println(err)
	}
}
