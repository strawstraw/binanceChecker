package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const (
	serverURL = "http://localhost:3001/api/v1/rates"
)

func GetCurrentRate(pair string) {
	url := fmt.Sprintf("%s?pairs=%s", serverURL, url.QueryEscape(pair))
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	rates := map[string]float64{}
	if err = json.Unmarshal(body, &rates); err != nil {
		log.Fatalln(err)
	}

	for _, rate := range rates {
		fmt.Println(rate)
	}

}
