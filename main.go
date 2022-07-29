package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"goForward/config"
	"io"
	"io/ioutil"
	"net/http"
)

func main() {
	config.ReadConfig()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/raw", Raw)
	router.HandleFunc("/comtroler", Comtroler)
	addr := fmt.Sprintf(":%d", config.Config.Port)
	log.Fatal(http.ListenAndServe(addr, router)) // addr is in form :9001
}

func Raw(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		log.Errorf("Raw: Error during HTTP receive %s", err)
		return
	}

	for _, url := range config.Config.Webhooks {
		go forwardWebhook(url, body, r.Header)
	}

	//if decodeError == true {
	//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//	w.WriteHeader(http.StatusUnprocessableEntity)
	//	return
	//}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func Comtroler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		log.Errorf("Raw: Error during HTTP receive %s", err)
		return
	}

	for _, url := range config.Config.Webhooks {
		go forwardWebhook(url, body, r.Header)
	}

	//if decodeError == true {
	//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//	w.WriteHeader(http.StatusUnprocessableEntity)
	//	return
	//}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func forwardWebhook(url string, body []byte, field http.Header) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))

	if err != nil {
		log.Warnf("Sender: unable to connect to %s - %s", url, err)
		return
	}

	req.Header.Set("X-Goforward", "Forwarded!")
	for k, v := range field {
		for _, x := range v {
			req.Header.Add(k, x)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Warningf("Webhook: %s", err)
		return
	}
	defer resp.Body.Close()

	log.Debugf("Webhook: Response %s", resp.Status)
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))
}
