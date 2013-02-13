package main

import (
    "net/http"
    "log"
    "github.com/bpina/go-tracker/thp"
    "github.com/bpina/go-tracker/tools"
    "github.com/bpina/go-tracker/data"
    "github.com/bpina/go-tracker/data/configuration"
)

var DbConfig configuration.DatabaseConfiguration

func AnnounceHandler(w http.ResponseWriter, req *http.Request) {
    log.Printf(req.URL.RawQuery)
    w.Header().Set("Content-Type", "text/plain")

    if req.Method == "POST" {
        response := thp.NewErrorResponse("Unsupported HTTP method.")
        w.Write([]byte(response.String()))
        log.Printf(response.String())
        return
    }

    err := data.OpenDatabaseConnection(DbConfig)
    if err != nil {
        response := thp.NewErrorResponse("No database connection.")
        w.Write([]byte(response.String()))
        log.Printf(response.String())
        return
    }
    defer data.CloseDatabaseConnection()

    req.ParseForm()
    announce, errors := thp.NewAnnounce(req.Form)
    log.Printf("something new")

    if errors != nil {
        message := tools.FormatErrors(errors)
        response := thp.NewErrorResponse(message)
        w.Write([]byte(response.String()))
        log.Printf(response.String())
        return
    } else {
        torrent, err := data.FindTorrent(announce.InfoHash)
        if err != nil {
            log.Printf(err.Error())
            response := thp.NewErrorResponse("Database error.")
            w.Write([]byte(response.String()))
            log.Printf(response.String())
            return
        }

        if torrent == nil {
            response := thp.NewErrorResponse("Could not locate torrent.")
            w.Write([]byte(response.String()))
            log.Printf(response.String())
            return
        }

        response := new(thp.Response)
        response.Interval = 30
        response.Complete = 999
        response.Incomplete = 999

        peer := new(thp.ConnectedPeer)
        peer.Ip = "192.168.3.99"
        peer.Port = 1337
        peer.PeerId = announce.PeerId

        response.Peers = append(response.Peers, *peer)

        message := response.String()
        log.Printf(message)
        w.Write([]byte(message))
       return
    }
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
