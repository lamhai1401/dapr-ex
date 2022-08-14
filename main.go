package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/gorilla/mux"
)

type stateData struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/checkout", getCheckout).Methods("POST", "OPTIONS")

	r.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("Echo body", string(body))
		}

		w.Write(body)
	})

	r.HandleFunc("/greeting",
		func(w http.ResponseWriter, r *http.Request) {
			resp, _ := http.Get("http://localhost:8089/v1.0/state/statestore/mystate")
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)

			strVal := string(body)
			count := 0
			if strVal != "" {
				count, _ = strconv.Atoi(strVal)
				count++
			}

			stateObj := []stateData{}
			stateObj = append(stateObj, stateData{Key: "mystate", Value: count})
			stateData, _ := json.Marshal(stateObj)

			resp, _ = http.Post("http://localhost:8089/v1.0/state/statestore", "application/json", bytes.NewBuffer(stateData))
			if count == 1 {
				fmt.Fprintf(w, "I’ve greeted you "+strconv.Itoa(count)+" time.")
			} else {
				fmt.Fprintf(w, "I’ve greeted you "+strconv.Itoa(count)+" times.")
			}
		},
	)

	go http.ListenAndServe(":8088", r)

	OutputBinding()
	// RunGRPCClient()
	select {}
}

//code
func getCheckout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var orderId int
	err := json.NewDecoder(r.Body).Decode(&orderId)
	log.Println("Received Message: ", orderId)
	if err != nil {
		log.Printf("error parsing checkout input binding payload: %s", err)
		w.WriteHeader(http.StatusOK)
		return
	}
}

func OutputBinding() {
	BINDING_NAME := "checkout"
	BINDING_OPERATION := "create"
	client, err := dapr.NewClient()
	if err != nil {
		panic(err.Error())
	}
	defer client.Close()

	for {
		time.Sleep(5000 * time.Millisecond)
		ctx := context.Background()
		orderId := rand.Intn(1000-1) + 1
		// Using Dapr SDK to invoke output binding
		in := &dapr.InvokeBindingRequest{Name: BINDING_NAME, Operation: BINDING_OPERATION, Data: []byte(strconv.Itoa(orderId))}
		err = client.InvokeOutputBinding(ctx, in)
		if err != nil {
			panic(err.Error())
		}
		log.Println("Sending message: " + strconv.Itoa(orderId))
	}
}
