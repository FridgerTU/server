package main

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
