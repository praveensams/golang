// Copyright 2018 Whatsapp, Inc. and its affiliates.
/// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func sendHttpsRequest(url string, headerKey string, headerVal string, requestType string, httpTimeout int) ([]byte, int) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	netClient := &http.Client{
		Timeout:   time.Second * time.Duration(httpTimeout),
		Transport: transport,
	}
	request, err := http.NewRequest(requestType, url, nil)
	if err != nil {
		log.Fatalln("Not able to create new request to", url, ":", err)
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

func dumpToFile(output string, dest string) {
	tmpfile, err := ioutil.TempFile("", "tmp")
	if err != nil {
		log.Println("Failed to create temp file:", err)
		return
	}
	defer os.Remove(tmpfile.Name()) // clean up
	content := []byte(output)
	if _, err := tmpfile.Write(content); err != nil {
		log.Println("Failed to write content to file", tmpfile.Name(), err)
		return
	}
	if err := tmpfile.Close(); err != nil {
		log.Println("Failed to close file", tmpfile.Name(), err)
		return
	}
	if err := os.Rename(tmpfile.Name(), dest); err != nil {
		log.Println("Failed to rename file", tmpfile.Name(), "to", dest, err)
		return
	}
	log.Println("Successfully written to file", dest)
}

func genTargetsFile(body []byte, url string, nodeExporterPort int, cadvisorPort int, targetsFile string) {
	var jsonResp map[string]interface{}
	if err := json.Unmarshal(body, &jsonResp); err != nil {
		log.Fatalln("Request ", url, "returns non-json response:", string(body))
	}
	health, exist := jsonResp["health"]
	if !exist || reflect.TypeOf(health).String() != "map[string]interface {}" {
		log.Fatalln("Request to", url, "returns bad response", string(body))
	}
	healthMap := health.(map[string]interface{})
	_, exist = healthMap["gateway_status"]
	var machines map[string]bool
	machines = make(map[string]bool)
	if exist {
		// single-connect
		coreEndpoint := os.Getenv("WA_CORE_ENDPOINT")
		machines[coreEndpoint] = true
	} else {
		// multi-connect
		for k := range healthMap {
			index := strings.LastIndex(k, ":")
			machines[k[0:index]] = true
		}
	}
	nodeExporters := make([]string, 0, len(machines))
	cadvisors := make([]string, 0, len(machines))
	for k := range machines {
		nodeExporters = append(nodeExporters, "\""+k+":"+strconv.Itoa(nodeExporterPort)+"\"")
		cadvisors = append(cadvisors, "\""+k+":"+strconv.Itoa(cadvisorPort)+"\"")
	}
	output := "["
	output += "{ \"targets\": [" + strings.Join(nodeExporters, ", ") + "], \"labels\" : { \"job\": \"node-exporter\"}},"
	output += "{ \"targets\": [" + strings.Join(cadvisors, ", ") + "], \"labels\" : { \"job\": \"cadvisor\"}}"
	output += "]"
	log.Println("Targest config:", output)
	dumpToFile(output, targetsFile)
}

func genAuthToken(url string, tokenFile string, httpTimeout int) string {
	username := os.Getenv("WA_WEB_USERNAME")
	password := os.Getenv("WA_WEB_PASSWORD")
	if username == "" {
		log.Fatalln("env WA_WEB_USERNAME is empty")
	}
	if password == "" {
		log.Fatalln("env WA_WEB_PASSWORD is empty")
	}
	base64Str := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	body, code := sendHttpsRequest(url, "Authorization", "Basic "+base64Str, "POST", httpTimeout)
	if code != 200 {
		log.Fatalln("Failed to get auth token from", url, "code:", code, "response:", string(body))
	}
	var jsonResp map[string]interface{}
	if err := json.Unmarshal(body, &jsonResp); err != nil {
		log.Fatalln("Request ", url, "returns non-json response:", string(body))
	}
	users, exist := jsonResp["users"]
	if !exist || reflect.TypeOf(users).String() != "[]interface {}" {
		log.Fatalln("Request to", url, "returns bad response:", string(body))
	}
	usersArray := users.([]interface{})
	if len(usersArray) == 0 || reflect.TypeOf(usersArray[0]).String() != "map[string]interface {}" {
		log.Fatalln("Request to", url, "returns bad users response:", string(body))
	}
	tokenMap := usersArray[0].(map[string]interface{})
	token, exist := tokenMap["token"]
	if !exist || reflect.TypeOf(token).String() != "string" {
		log.Fatalln("Request to", url, "returns bad response:", string(body))
	}
	tokenStr := token.(string)
	log.Println("Got bearer token:", tokenStr)
	dumpToFile(tokenStr, tokenFile)
	return tokenStr
}

func main() {
	webEndpoint := os.Getenv("WA_WEB_ENDPOINT")
	if webEndpoint == "" {
		log.Fatalln("env WA_WEB_ENDPOINT is empty, exit!")
	}
	loginURL := "https://" + webEndpoint + "/v1/users/login"
	healthURL := "https://" + webEndpoint + "/v1/health"
	log.Println("web endpoint:", webEndpoint)
	tokenFile := flag.String("tokenfile", "/etc/prometheus/auth_token", "locaiton of file containing bearer token")
	httpTimeout := flag.Int("timeout", 30, "Timeout of https request in seconds")
	nodeExporterPort := flag.Int("node-exporter-port", 9100, "Port of node exporter container")
	cadvisorPort := flag.Int("cadvisor-port", 8080, "Port of cadvisor container")
	targetsFile := flag.String("targetsfile", "/etc/prometheus/targets.json", "location of file containing targets configs for node and container monitoring")
	flag.Parse()
	log.Println("token file:", *tokenFile)
	log.Println("http timeout:", *httpTimeout)
	log.Println("node exporter port:", *nodeExporterPort)
	log.Println("cadvisor port:", *cadvisorPort)
	log.Println("target file:", *targetsFile)
	token := ""
	if _, err := os.Stat(*tokenFile); err == nil {
		file, err := os.Open(*tokenFile)
		if err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				if strings.TrimSpace(scanner.Text()) != "" {
					token = strings.TrimSpace(scanner.Text())
					break
				}
			}
		}
	}
	if token == "" {
		token = genAuthToken(loginURL, *tokenFile, *httpTimeout)
	}
	healthBody, code := sendHttpsRequest(healthURL, "Authorization", "Bearer "+token, "GET", *httpTimeout)
	if code != 200 {
		token = genAuthToken(loginURL, *tokenFile, *httpTimeout)
	}
	genTargetsFile(healthBody, healthURL, *nodeExporterPort, *cadvisorPort, *targetsFile)
}
