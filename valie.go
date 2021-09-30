package main

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

type S struct {
	count int
}

func sendHttpsRequest(url string, headerKey string, headerVal string, requestType string, httpTimeout int) ([]byte, int) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	netClient := &http.Client{
		Timeout:   time.Second * time.Duration(httpTimeout),
		Transport: transport,
	}
	request, err := http.NewRequest(requestType, url+"/v1/users/login", nil)
	login := url + "/v1/users/login"
	if err != nil {
		log.Fatalln("Not able to create new request to", login, ":", err)
	}
	request.Header.Set(headerKey, headerVal)
	response, err := netClient.Do(request)
	if err != nil {
		log.Fatalln("Error to send request to", url, ":", err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln("Error to read response body from", url, ":", err)
	}
	return body, response.StatusCode
}

func (s *S) validate(url string, c chan string) {

	base64Str := base64.StdEncoding.EncodeToString([]byte("admin" + ":" + "oph>aX8abioM%ahrei5av2iro4we3Chi^ik+ae1x"))
	url1 := []string{}
	body, _ := sendHttpsRequest(url, "Authorization", "Basic "+base64Str, "POST", 10)
	r := regexp.MustCompile(`.*token\"\:\"(?P<token>[\w\d\.\-\_]+)`)
	lns := len(r.FindStringSubmatch(string(body)))
	appender := []string{"/v1/stats/app?format=prometheus", "/metrics?format=prometheus"}
	var form string
	for f := 0; f < len(appender); f++ {
		form = url + appender[f]
		url1 = append(url1, form)
	}

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + r.FindStringSubmatch(string(body))[lns-1]

	// Create a new request using http

	for i := 0; i < len(url1); i++ {

		req, err := http.NewRequest("GET", url1[i], nil)

		// add authorization header to the req
		req.Header.Add("Authorization", bearer)

		// Send req using http Client
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error on response.\n[ERROR] -", err)
		}
		defer resp.Body.Close()

		body1, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error while reading the response bytes:", err)
		}
		s.count = s.count + 1
		c <- string([]byte(body1))

	}

}

func page(w http.ResponseWriter, r *http.Request) {
	c := make(chan string, 100)
	count := &S{0}

	url1 := "https://prod-tier-uquuz3-wa.wab.cloud.unifonic.com"

	go count.validate(url1, c)

	time.Sleep(4 * time.Second)
	for j := 0; j < count.count; j++ {
		fmt.Fprintf(w, <-c)
	}
	fmt.Println("Completed")
}

func calls() {
	http.HandleFunc("/", page)
	log.Fatal(http.ListenAndServe(":9101", nil))

}

func main() {
	calls()
}
