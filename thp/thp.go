package thp

import (
    "github.com/bpina/go-tracker/data"
    "github.com/zeebo/bencode"
)

type ConnectedPeer struct {
    PeerId      string `bencode:"peer id"`
    Ip          string `bencode:"ip"`
    Port        int `bencode:"port"`
}


type Response struct {
    FailureReason   string
    Interval        int
    Complete        int
    Incomplete      int
    Peers           []ConnectedPeer `bencode:"peers"`
}

func NewErrorResponse(message string) *Response {
    return &Response{FailureReason: message}
}

func NewDatabaseErrorResponse() *Response {
    return NewErrorResponse("Database error.")
}

func NewTorrentResponse(torrent *data.Torrent, peers []data.Peer) *Response {
  response := new(Response)
  response.Interval = 30
  response.Complete = torrent.Complete
  response.Incomplete = torrent.Incomplete

  for i := range peers {
    connectedPeer := new(ConnectedPeer)
    connectedPeer.Ip = peers[i].Ip
    connectedPeer.Port = peers[i].Port
    connectedPeer.PeerId = peers[i].PeerId
    response.Peers = append(response.Peers, *connectedPeer)
  }

  return response
}

func (r Response) String() string {
    //TODO: make this a whole lot better
    if r.FailureReason != "" {
        var failure struct {
            FailureReason string `bencode:"failure reason"`
        }

        failure.FailureReason = r.FailureReason

        data, err := bencode.EncodeString(failure)
        if err != nil {
            return "failure encoding response"
        }

        return data
    }

    var response struct {
        Interval        int `bencode:"interval"`
        Complete        int `bencode:"complete"`
        Incomplete      int `bencode:"incomplete"`
        Peers           []ConnectedPeer `bencode:"peers"`
    }

    response.Interval = r.Interval
    response.Complete = r.Complete
    response.Incomplete = r.Incomplete
    response.Peers = r.Peers

    data, err := bencode.EncodeString(response)
    if err != nil {
        return "failure encoding response"
    }

    return data
}
