package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Response struct {
	title   string   `json:"title"`
	version string   `json:"version"`
	href    string   `json:"href"`
	Recipes []Recipe `json:"results"`
}

type Recipe struct {
	Title       string `json:"title"`
	href        string `json:"href"`
	Ingredients string `json:"ingredients"`
	Thumbnail   string `json:"thumbnail"`
}

const testUrl = "http://www.recipepuppy.com/api/?i=onions,garlic&q=omelet&p=3"

func test() {
	bytes, err := executeGetRequest(testUrl)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp := Response{}
	err = json.Unmarshal(bytes, &resp)
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