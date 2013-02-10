package main

import (
    "net/http"
    "log"
    "github.com/bpina/go-tracker/tracker"
    "github.com/bpina/go-tracker/tools"
)

func AnnounceHandler(w http.ResponseWriter, req *http.Request) {
    if req.Method == "POST" {
        response := tracker.NewErrorResponse("Unsupported HTTP method.")
        w.Write([]byte(response.String()))
        return
    }

    req.ParseForm()
    announce, err := tracker.NewAnnounce(req.Form)
    if err != nil {
        message := tools.FormatErrors(err)
        response := tracker.NewErrorResponse(message)
        w.Write([]byte(response.String()))
        return
    } else {
       announce.NumWant = 5
       return
    }
}

func main() {
    http.HandleFunc("/announce", AnnounceHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
