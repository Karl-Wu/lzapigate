package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	transport = &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 10 * time.Second,
	}

	authToken string
)

type gqlDataStruct map[string]interface{}

type gqlRespStruct struct {
	Data gqlDataStruct `json:"data"`
}

func postGraphQL(url, queryStr string) (gqlDataStruct, error) {
	query := map[string]string{
		"query": queryStr,
	}
	queryBytes, _ := json.Marshal(query)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(queryBytes))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// We need to set the content type from the writer, it includes necessary boundary as well
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

	// Do the request
	client := &http.Client{}
	resp, err := client.Do(req)

	// DEBUG.Println(GQL, "Post graphql body", resp)
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// DEBUG.Println(GQL, "Post graphql resp body", string(respBody))

	var gqlResp gqlRespStruct
	err = json.Unmarshal(respBody, &gqlResp)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	DEBUG.Println(GQL, "graphql data:", gqlResp.Data)

	return gqlResp.Data, nil
}
