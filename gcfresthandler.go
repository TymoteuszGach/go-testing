package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type GCFRestHandler struct {
	dbAdapter DatabaseAdapter
}

func NewGCFRestHandler(dbAdapter DatabaseAdapter) *GCFRestHandler {
	return &GCFRestHandler{dbAdapter: dbAdapter}
}

func (handler *GCFRestHandler) ListGCFResults(w http.ResponseWriter, r *http.Request) {
	results, err := handler.dbAdapter.ListGCFResults()
	resultsStrings := make([]string, len(results))
	for i, result := range results{
		resultsStrings[i] = result.String()
	}
	if err != nil {
		log.Printf("cannot list GCF results: %v\n", err)
		// set error response
		return
	}

	_, err = fmt.Fprintf(w, "Results:\n%s", strings.Join(resultsStrings, "\n"))
	if err != nil {
		log.Printf("cannot return GCF results: %v\n", err)
		return
	}
}