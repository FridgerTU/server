package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	controller := NewController()
	if err := http.ListenAndServe(fmt.Sprintf(":%v", getPort()), controller); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func getPort() string {
	if configuredPort := os.Getenv("PORT"); configuredPort == "" {
		return "8080"
	} else {
		return configuredPort
	}
}
