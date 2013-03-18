package main

import (
    "net"
    "net/http"
    "log"
    "github.com/bpina/go-tracker/thp"
    "github.com/bpina/go-tracker/tools"
    "github.com/bpina/go-tracker/data"
    "github.com/bpina/go-tracker/data/configuration"
)

var DbConfig configuration.DatabaseConfiguration

func WriteErrorResponse(w http.ResponseWriter, message string) {
  response := thp.NewErrorResponse(message)
  w.Write([]byte(response.String()))
}

func AnnounceHandler(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/plain")

    if req.Method != "GET" {
        response := thp.NewErrorResponse("Unsupported HTTP method.")
        w.Write([]byte(response.String()))
        return
    }

    err := data.OpenDatabaseConnection(DbConfig)
    if err != nil {
        response := thp.NewErrorResponse("No database connection.")
        w.Write([]byte(response.String()))
        return
    }
    defer data.CloseDatabaseConnection()

    req.ParseForm()
    announce, errors := thp.NewAnnounce(req.Form)
    if errors != nil {
        message := tools.FormatErrors(errors)
        response := thp.NewErrorResponse(message)
        w.Write([]byte(response.String()))
        return
    } else {
        host, _, err := net.SplitHostPort(req.RemoteAddr)
        if err != nil {
          WriteErrorResponse(w, "Could not determine remote host")
        }

        ip := net.ParseIP(host)
        if ip == nil {
          WriteErrorResponse(w, "Could not determine remote host")
        }

        is_ipv6 := ip.To4() == nil

        torrent, err := data.FindTorrent(announce.InfoHash)
        if err != nil {
            response := thp.NewErrorResponse("Database error.")
            w.Write([]byte(response.String()))
            return
        }

        if torrent == nil {
            response := thp.NewErrorResponse("Could not locate torrent.")
            w.Write([]byte(response.String()))
            return
        }

        peer, err := data.FindPeerByPeerIdAndInfoHash(announce.PeerId, torrent.InfoHash)
        if err != nil {
            response := thp.NewErrorResponse("Database error")
            w.Write([]byte(response.String()))
            return
        }

        if peer == nil {
            peer = new(data.Peer)
            peer.PeerId = announce.PeerId
            peer.Ip = host
            peer.Port = announce.Port
            peer.InfoHash = torrent.InfoHash
            peer.IsIpV6 = is_ipv6
            peer.Save()

            if announce.NumWant == 0 {
                torrent.Complete += 1
            } else {
                torrent.Incomplete += 1
            }
            torrent.Update()
        } else {
          peer.Ip = host
          peer.IsIpV6 = is_ipv6
          peer.Update()
        }

        response := new(thp.Response)
        response.Interval = 30
        response.Complete = torrent.Complete
        response.Incomplete = torrent.Incomplete

        availablePeers, err := data.FindAvailablePeers(announce.PeerId, torrent.InfoHash, is_ipv6)
        if err != nil {
          WriteErrorResponse(w, "Database error")
          return
        }

        for i := range availablePeers {
          connectedPeer := new(thp.ConnectedPeer)
          connectedPeer.Ip = availablePeers[i].Ip
          connectedPeer.Port = availablePeers[i].Port
          connectedPeer.PeerId = availablePeers[i].PeerId
          response.Peers = append(response.Peers, *connectedPeer)
        }

        message := response.String()
        w.Write([]byte(message))
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
