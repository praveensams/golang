package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var wg sync.Mutex
var vge sync.WaitGroup

type S struct {
	x int
}

type SlackRequestBody struct {
	Text string `json:"text"`
}

func (count *S) validate(i int, se string, c chan string) {
	h := se
	addr := fmt.Sprintf("%s:%d", h, i)
	s, err := net.DialTimeout("tcp", addr, time.Second*2)
	if err == nil {
		count.x = count.x + 1
		ads := fmt.Sprintf("%s  -> %d", se, i)
		c <- ads
		defer s.Close()

	}

	defer vge.Done()
}

func main() {
	file, err := os.Open("url.txt")
	if err != nil {
		log.Fatalf("failed to open")
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	file.Close()

	count := &S{
		x: 0,
	}
	c := make(chan string, 10)
	count.x = 0
	if len(text) == 0 {
		panic("File is empty , Please verify the content")
	}
	for _, j := range text {
		for i := 0; i < 1024; i++ {
			vge.Add(1)
			go count.validate(i, j, c)
		}
	}

	vge.Wait()
	time.Sleep(4 * time.Second)
	webhookUrl := "https://hooks.slack.com/services/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"
	t0 := time.Now().String()[:25]
	full := fmt.Sprintf("Scanning Results  => %s", t0)
	err = SendSlackNotification(webhookUrl, full)
	if err != nil {
		log.Fatal(err)
	}
	lists := make([]string, 100)
	for j := 0; j < count.x; j++ {
		lists = append(lists, <-c)
	}
	err = SendSlackNotification(webhookUrl, strings.Join(lists, "\n"))
	if err != nil {
		log.Fatal(err)

	}
}

func SendSlackNotification(webhookUrl string, msg string) error {

	slackBody, _ := json.Marshal(SlackRequestBody{Text: msg})
	req, err := http.NewRequest(http.MethodPost, webhookUrl, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("Non-ok response returned from Slack")
	}
	return nil
}
