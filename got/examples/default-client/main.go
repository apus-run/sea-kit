package main

import (
	"fmt"
	"net/http"

	"github.com/apus-run/sea-kit/http_client"
)

func main() {
	// Create a new client with default options.
	client := got.NewClient()

	// Create a new request with default options.
	req, err := http.NewRequest(http.MethodGet, "https://httpbin.org/anything", nil)

	// Send the request and get the response.
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// Hydrate the response body into a safemap.
	var target map[string]any
	err = got.ResponseToJson(resp, &target)

	// Print the target safemap.
	fmt.Println(target)
}
