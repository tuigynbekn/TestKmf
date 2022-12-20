package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type ProxyRequest struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
}

type ProxyResponse struct {
	ID      string              `json:"id"`
	Status  int                 `json:"status"`
	Length  int                 `json:"length"`
	Headers map[string][]string `json:"headers"`
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	proxyReq := ProxyRequest{}
	err := json.NewDecoder(r.Body).Decode(&proxyReq)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest(proxyReq.Method, proxyReq.URL, nil)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for key, value := range proxyReq.Headers {
		req.Header.Add(key, value)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	proxyRes := ProxyResponse{
		ID:      uuid.New().String(),
		Status:  res.StatusCode,
		Length:  len(body),
		Headers: map[string][]string{},
	}

	for key, values := range r.Header {
		proxyRes.Headers[key] = values
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response, err := json.Marshal(proxyRes)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(response)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/proxy", proxyHandler)
	http.Handle("/", router)

	fmt.Println("Server is listening... http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
