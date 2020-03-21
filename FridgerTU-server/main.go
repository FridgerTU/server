package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

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
}
