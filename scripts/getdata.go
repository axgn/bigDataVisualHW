package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func posturl(client *http.Client, url string, paras map[string]string) string {

	postbody, err := json.Marshal(paras)
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Post(url, "application/json", bytes.NewBuffer(postbody))
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}

func geturl(client *http.Client, url string, paras map[string]string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	q := req.URL.Query()
	for k, v := range paras {
		q.Set(k, v)
	}
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}
func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	var client http.Client
	body1 := geturl(&client, "http://www.google.com/robots.txt", nil)
	fmt.Printf("%s\n\n", body1[:100])
	body2 := posturl(&client, "https://httpbin.org/post", nil)
	fmt.Printf("%s\n\n", body2[:100])
}
