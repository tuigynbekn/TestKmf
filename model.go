package main

import "sync"

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

type History struct {
	ProxyRequest  `json:"request"`
	ProxyResponse `json:"response"`
}

type HistoryMap struct {
	sync.RWMutex
	History map[string]History `json:"history"`
}

func NewHistoryMap() *HistoryMap {
	return &HistoryMap{
		History: make(map[string]History),
	}
}

func (rm *HistoryMap) GetHistory(key string) (value History) {
	rm.RLock()
	value = rm.History[key]
	rm.RUnlock()
	return
}

func (rm *HistoryMap) SaveHistory(key string, value History) {
	rm.Lock()
	rm.History[key] = value
	rm.Unlock()
}
