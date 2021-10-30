package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	fmt.Println(os.Getenv("HTTP_PROXY"))
	// u, err := url.Parse("http://localhost:8000")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	client := &http.Client{
		// Transport: &http.Transport{Proxy: http.ProxyURL(u)},
	}
	resp, err := client.Get("http://localhost:8090/vm/1")
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}
