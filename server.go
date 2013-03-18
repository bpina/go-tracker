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
  log.Printf(response.String())
}

func AnnounceHandler(w http.ResponseWriter, req *http.Request) {
    log.Printf(req.RemoteAddr)
    w.Header().Set("Content-Type", "text/plain")

    if req.Method != "GET" {
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
    if errors != nil {
        message := tools.FormatErrors(errors)
        response := thp.NewErrorResponse(message)
        w.Write([]byte(response.String()))
        log.Printf(response.String())
        return
    } else {
        host, port, err := net.SplitHostPort(req.RemoteAddr)
        if err != nil {
          WriteErrorResponse(w, "Could not determine remote host")
        }
        log.Printf(port)

        ip := net.ParseIP(host)
        if ip == nil {
          WriteErrorResponse(w, "Could not determine remote host")
        }

        is_ipv6 := ip.To4() == nil

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
            log.Printf(err.Error())
            response := thp.NewErrorResponse("Database error")
            w.Write([]byte(response.String()))
            log.Printf(response.String())
            return
        }

        if peer == nil {
            log.Printf("Making a new peer")
            peer = new(data.Peer)
            peer.PeerId = announce.PeerId
            peer.Ip = host
            peer.Port = announce.Port
            peer.InfoHash = torrent.InfoHash
            peer.IsIpV6 = is_ipv6
            err = peer.Save()
            if err != nil {
              log.Printf(err.Error())
            }

            if announce.NumWant == 0 {
                torrent.Complete += 1
            } else {
                torrent.Incomplete += 1
            }
            err = torrent.Update()
            if err != nil {
              log.Printf(err.Error())
            }
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
          log.Printf(err.Error())
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
        log.Printf(message)
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
