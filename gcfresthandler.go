package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type ListGCFResults func() ([]GCFResult, error)

type GCFRestHandler struct {
	listGCFResults ListGCFResults
}

func NewGCFRestHandler(listGCFResults ListGCFResults) *GCFRestHandler {
	return &GCFRestHandler{listGCFResults: listGCFResults}
}

func (handler *GCFRestHandler) ListGCFResults(w http.ResponseWriter, r *http.Request) {
	results, err := handler.listGCFResults()
	resultsStrings := make([]string, len(results))
	for i, result := range results {
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
