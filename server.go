package main

import (
    "net/http"
    "log"
    "go-tracker/tracker"
)

func AnnounceHandler(w http.ResponseWriter, req *http.Request) {
    if req.Method == "POST" {
        return tracker.NewErrorResponse("Unsupported")
    }

    req.ParseForm()
    log.Print(req.Form.Encode())

    announce := tracker.NewAnnounce(req.Form)
    if announce == nil {
        log.Print("announce was bad")
    } else {
        log.Print("announce was good")
    }
}

func main() {
    http.HandleFunc("/announce", AnnounceHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}