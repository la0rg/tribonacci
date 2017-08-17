package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/la0rg/tribonacci/tribo"
)

// TribHandler is a base struct to keep all dependencies
// and http handlers of the programm.
type TribHandler struct {
	tribo *tribo.Tribo
}

// NewTribHandler creates new instance of TribHandler.
func NewTribHandler() *TribHandler {
	return &TribHandler{
		tribo: tribo.New(100000),
	}
}

// TribonacciHandler handles requests for getting tribonacci sequence value by its serial number.
func (t *TribHandler) TribonacciHandler(w http.ResponseWriter, r *http.Request) {
	n, ok := mux.Vars(r)["n"]
	if !ok {
		http.Error(w, "N should be specified", http.StatusBadRequest)
		return
	}
	number, err := strconv.Atoi(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx, cn := context.WithTimeout(r.Context(), time.Minute)
	defer cn()
	result, err := t.tribo.Get(ctx, number)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
