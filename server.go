package main

import (
    "net/http"
    "log"
    "github.com/bpina/go-tracker/thp"
    "github.com/bpina/go-tracker/data/configuration"
)

var DbConfig configuration.DatabaseConfiguration

func AnnounceHandler(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    var response *thp.Response

    tracker, err := thp.NewTracker(DbConfig, req)
    if err != nil {
      response = thp.NewErrorResponse("Failed to initialize tracker.")
    } else {
      response = tracker.Execute()
    }

    w.Write([]byte(response.String()))
}

func main() {
    var err error
    DbConfig, err = configuration.NewDatabaseConfiguration()

    if err != nil {
        panic(err)
    }

    http.HandleFunc("/announce", AnnounceHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
