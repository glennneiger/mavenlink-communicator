package api

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"
)

// BasicAuth uses the provided user name and password to generate
// a base64 encoded string that can be utilised for basic
// authentication purposes
func BasicAuth(username, password string) string {
	// format username and password
	auth := username + ":" + password
	// encode to base64 string
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// InsecureRequest makes an HTTP call to the provided endpoint(param: url)
// using the prescribed HTTP request type(param: method). The response from
// the endpoint(param: url) is then decoded by the json package into the
// specified structure(param: target). The request is insecure due to the HTTP
// transport layer being configured with the {InsecureSkipVerify: true} option
func InsecureRequest(url string, method string, body interface{}, token string, target interface{}) error {
	// set the TLS option to skip insecure certificates
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	// create a HTTP client with a timeout of 30 secs
	client := &http.Client{Transport: tr, Timeout: time.Second * time.Duration(30)}

	// format JSON body
	var rawBody bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&rawBody).Encode(body); err != nil {
			return err
		}
	}

	// create a new HTTP request
	httpReq, requestErr := http.NewRequest(method, url, &rawBody)
	if requestErr != nil {
		log.Printf("HTTP Request Error : %s\n", requestErr)
		return nil
	}
	// add custom headers
	httpReq.Header.Add("Content-Type", "application/json")
	// add authentication user-agent to header
	httpReq.Header.Set("User-Agent", "mavenlink-communicator/1.0")
	// add authentication token to header
	httpReq.Header.Add("Authorization", "Bearer "+token)
	log.Printf(":::::::::: Mavenlink Communicator HTTP request ::::::::::")
	log.Printf("==>> %v <<==", httpReq)
	// use the HTTP client to perform the HTTP request
	httpResp, requestErr := client.Do(httpReq)
	// check request has no error
	if requestErr != nil {
		log.Printf("HTTP Request Error : %s\n", requestErr)
		return nil
	}
	// check response for status : NOT FOUND, UNAUTHORISED & FORBIDDEN
	if 404 == httpResp.StatusCode ||
		401 == httpResp.StatusCode ||
		403 == httpResp.StatusCode ||
		400 == httpResp.StatusCode ||
		500 == httpResp.StatusCode {
		formattedErr := fmt.Sprintf("Error : %s(Status: %d)", httpResp.Status, httpResp.StatusCode)
		return errors.New(formattedErr)
	}
	// check response for status : NO CONTENT
	if 204 == httpResp.StatusCode {
		return nil
	}
	// decode response body to the intended target structure
	return json.NewDecoder(httpResp.Body).Decode(target)
}
