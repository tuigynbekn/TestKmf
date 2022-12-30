package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	historyMap *HistoryMap
}

func NewHandler() *Handler {
	return &Handler{
		historyMap: NewHistoryMap(),
	}
}

func (h *Handler) proxyHandler(w http.ResponseWriter, r *http.Request) {
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

	history := History{
		ProxyRequest:  proxyReq,
		ProxyResponse: proxyRes,
	}
	h.historyMap.SaveHistory(proxyRes.ID, history)

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

func (h *Handler) historyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var data any
	switch mux.Vars(r)["id"] {
	case "":
		data = h.historyMap.History
	default:
		data = h.historyMap.GetHistory(mux.Vars(r)["id"])
	}

	res, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(res)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
