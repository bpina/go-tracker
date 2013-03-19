package thp

import (
    "net"
    "net/http"
    "github.com/bpina/go-tracker/data/configuration"
    "github.com/bpina/go-tracker/data"
    "github.com/bpina/go-tracker/tools"
)

type Tracker struct {
    RemoteHost  string
    IsIpV6      bool
    Request     *http.Request
    DbConfig    configuration.DatabaseConfiguration
}

func NewTracker(dbConfig configuration.DatabaseConfiguration, req *http.Request) (t *Tracker, err error) {
    err = ValidateTrackerRequest(req)
    if err != nil {
      return t, err
    }

    ip, err := ParseRemoteAddress(req.RemoteAddr)
    if err != nil {
      return t, err
    }

    t = new(Tracker)
    t.RemoteHost = string(ip)
    t.IsIpV6 = ip.To4() == nil
    t.Request = req
    t.DbConfig = dbConfig

    return t, err
}

func (t *Tracker) Execute() *Response {
  announce, errors := NewAnnounce(t.Request)
  if errors != nil {
    message := tools.FormatErrors(errors)
    return NewErrorResponse(message)
  }

  err := data.OpenDatabaseConnection(t.DbConfig)
  if err != nil {
      return NewDatabaseErrorResponse()
  }
  defer data.CloseDatabaseConnection()

  torrent, err := data.FindTorrent(announce.InfoHash)
  if err != nil {
      return NewDatabaseErrorResponse()
  }

  if torrent == nil {
      return NewErrorResponse("Could not locate torrent.")
  }

  torrent.Adjust(announce.NumWant)

  peer, err := data.FindPeerByPeerIdAndInfoHash(announce.PeerId, torrent.InfoHash)
  if err != nil {
      return NewDatabaseErrorResponse()
  }

  if peer == nil {
    peer = t.CreatePeer(announce)
  } else {
    t.UpdatePeer(peer, announce)
  }

  peers, err := torrent.GetPeers(peer)
  if err != nil {
    return NewDatabaseErrorResponse()
  }

  return NewTorrentResponse(torrent, peers)
}

func (t *Tracker) CreatePeer(announce *Announce) *data.Peer {
  peer := new(data.Peer)
  peer.PeerId = announce.PeerId
  peer.Ip = t.RemoteHost
  peer.Port = announce.Port
  peer.InfoHash = announce.InfoHash
  peer.IsIpV6 = t.IsIpV6
  peer.Save()

  return peer
}

func (t *Tracker) UpdatePeer(peer *data.Peer, announce *Announce) {
  peer.Port = announce.Port
  peer.Ip = t.RemoteHost
  peer.IsIpV6 = t.IsIpV6
  peer.Update()
}

func ValidateTrackerRequest(req *http.Request) error {
  if req.Method != "GET" {
    return &TrackerError{"Unsupported HTTP method"}
  }

  if req.RemoteAddr == "" {
    return &TrackerError{"Unable to identify remote address"}
  }

  return nil
}

func ParseRemoteAddress(remoteAddress string) (ip net.IP, err error) {
  host, _, err := net.SplitHostPort(remoteAddress)
  if err != nil {
    return ip, err
  }

  ip = net.ParseIP(host)
  if ip == nil {
    err = &TrackerError{"Could not parse remote address"}
  }
  return ip, err
}

type TrackerError struct {
  Message string
}

func (e *TrackerError) Error() string {
  return e.Message
}
