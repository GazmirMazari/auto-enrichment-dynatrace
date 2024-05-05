package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	TraceparentHeader = "traceparent"
)

// APIResponse represents a simplified structure for the JSON response
type APIResponse struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// AddTraceHeader adds W3C Trace Context headers
func AddTraceHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a unique trace ID
		traceID := uuid.New()
		traceparent := fmt.Sprintf("00-%s-%s-01", traceID, uuid.New())

		// Add the W3C trace context header
		r = r.WithContext(context.WithValue(r.Context(), TraceparentHeader, traceparent))
		w.Header().Set(TraceparentHeader, traceparent)

		next.ServeHTTP(w, r)
	})
}

// FetchDataFromAPI retrieves data from a public API
func FetchDataFromAPI() ([]APIResponse, error) {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/todos")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse []APIResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, nil
}

// ProcessData simulates data processing and computation
func ProcessData(data []APIResponse) []APIResponse {
	var processedData []APIResponse
	for _, item := range data {
		// Simulate some data processing
		if item.Completed {
			item.Title = fmt.Sprintf("Processed: %s", item.Title)
			processedData = append(processedData, item)
		}
	}
	time.Sleep(500 * time.Millisecond) // Simulate further processing delay
	return processedData
}

// FetchAndProcessHandler handles the data fetching and processing
func FetchAndProcessHandler(w http.ResponseWriter, r *http.Request) {
	data, err := FetchDataFromAPI()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	processedData := ProcessData(data)
	response, err := json.Marshal(processedData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func main() {
	r := mux.NewRouter()
	r.Use(AddTraceHeader)
	r.Use(TraceMiddleware)

	// Define a simple handler that responds to requests
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	}).Methods("GET")

	// Add the FetchAndProcessHandler route
	r.HandleFunc("/fetch-and-process", FetchAndProcessHandler).Methods("GET")

	// Start theclear
	//server
	go http.ListenAndServe(":8080", r)

	// Simulate constant traffic
	endTime := time.Now().Add(1 * time.Minute)
	for time.Now().Before(endTime) {
		_, _ = http.Get("http://localhost:8080/fetch-and-process")
		time.Sleep(2 * time.Second) // Adjust the frequency as needed
	}
}
