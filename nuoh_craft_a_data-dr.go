package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Notification struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type Config struct {
	APIKey      string   `json:"api_key"`
	APIEndpoint string   `json:"api_endpoint"`
	Interval    int      `json:"interval"`
	Notifications []Notification `json:"notifications"`
}

func main() {
	var configFile string
	flag.StringVar(&configFile, "c", "config.json", "config file")
	flag.Parse()

 configFileBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
	log.Fatal(err)
 }

 var config Config
 err = json.Unmarshal(configFileBytes, &config)
 if err != nil {
	log.Fatal(err)
 }

 ticker := time.NewTicker(time.Duration(config.Interval) * time.Second)
 quit := make(chan struct{})
 go func() {
	for {
		select {
		case <-ticker.C:
			sendNotifications(config)
		case <-quit:
			ticker.Stop()
			return
		}
	}
}()

 fmt.Println("Notifier started. Press Ctrl+C to quit.")
 <-quit
}

func sendNotifications(config Config) {
	client := &http.Client{}
	for _, notification := range config.Notifications {
		req, err := http.NewRequest("POST", config.APIEndpoint, nil)
		if err != nil {
			log.Println(err)
			continue
		}
		req.Header.Set("Authorization", "Bearer "+config.APIKey)
		req.Header.Set("Content-Type", "application/json")
		jsonData, _ := json.Marshal(notification)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(jsonData))

		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Println("Error sending notification:", resp.Status)
		}
	}
}