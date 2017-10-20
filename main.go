package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Message struct {
	Text       string `json:"text"`
	IsMarkdown bool   `json:"mrkdwn"`
}

func main() {
	var text string
	flag.StringVar(&text, "t", "", "message text")

	var fpath string
	flag.StringVar(&fpath, "f", "", "message file path")

	flag.Parse()

	if len(text) == 0 && len(fpath) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	url := os.Getenv("SLACK_URL")
	if len(url) == 0 {
		log.Fatalf("[ERROR] SLACK_URL must be set in environment")
	}

	msg := Message{}
	if len(text) > 0 {
		msg.Text = text
	} else {
		b, err := ioutil.ReadFile(fpath)
		if err != nil {
			log.Fatalf("[ERROR] file read failed: %v", err)
		}
		msg.Text = string(b)
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("[ERROR] json marshal failed: %v", err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		log.Fatalf("[ERROR] slack post failed: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("[INFO] success")
}
