package main

import (
	"github.com/bpina/go-tracker/data"
	"github.com/bpina/go-tracker/data/configuration"
	"github.com/bpina/go-tracker/thp"
	"log"
	"net/http"
)

func AnnounceHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	var response *thp.Response

	tracker, err := thp.NewTracker(req)
	if err != nil {
		response = thp.NewErrorResponse("Failed to initialize tracker.")
	} else {
		response = tracker.Execute()
	}

	w.Write([]byte(response.String()))
}

func main() {
	dbConfig, err := configuration.NewDatabaseConfiguration()
	if err != nil {
		panic(err)
	}

	data.Database, err = data.OpenDatabaseConnection(dbConfig)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/announce", AnnounceHandler)
	log.Fatal(http.ListenAndServe(":9000", nil))
}
