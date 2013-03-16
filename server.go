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
            newTorrent := new(data.Torrent)
            newTorrent.InfoHash = announce.InfoHash
            newTorrent.Complete = 0
            newTorrent.Incomplete = 0

            err = newTorrent.Save()
            if err != nil {
                response := thp.NewErrorResponse(err.Error())
                w.Write([]byte(response.String()))
                log.Printf(response.String())
                return
            }

            response := thp.NewErrorResponse("Could not locate torrent.")
            w.Write([]byte(response.String()))
            log.Printf(response.String())
            return
        }

        peer, err := data.FindPeerByPeerIdAndInfoHash(announce.PeerId, torrent.InfoHash)
        if err != nil {
            response := thp.NewErrorResponse("Database error")
            w.Write([]byte(response.String()))
            log.Printf(response.String())
        }

        if peer == nil {
            peer = new(data.Peer)
            peer.PeerId = announce.PeerId
            peer.Ip = announce.Ip
            peer.Port = announce.Port
            peer.InfoHash = torrent.InfoHash
            peer.Save()

            if announce.NumWant == 0 {
                torrent.Complete += 1
            } else {
                torrent.Incomplete += 1
            }
            torrent.Update()
        }

        if announce.IpV6 != "" {
            peer.Ip = announce.Ip
            peer.IsIpV6 = false
            peer.Update()
        } else {
            peer.Ip = announce.IpV6
            peer.IsIpV6 = true
            peer.Update()
        }

        response := new(thp.Response)
        response.Interval = 30
        response.Complete = torrent.Complete
        response.Incomplete = torrent.Incomplete

        connectedPeer := new(thp.ConnectedPeer)
        connectedPeer.Ip = "192.168.3.99"
        connectedPeer.Port = 1337
        connectedPeer.PeerId = announce.PeerId

        response.Peers = append(response.Peers, *connectedPeer)

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
