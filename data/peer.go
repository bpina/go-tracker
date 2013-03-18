package data

import (
    _ "github.com/bmizerany/pq"
    "database/sql"
    "strconv"
)

type Peer struct {
    Id          int
    PeerId      string
    Ip          string
    Port        int
    InfoHash    string
    IsIpV6      bool
}

func (p *Peer) GetFields() map[string] string {
    fields := map[string] string {
        "peer_id": "'" + p.PeerId + "'",
        "ip": "'" + p.Ip + "'",
        "port": strconv.Itoa(p.Port),
        "info_hash": "'" + p.InfoHash + "'",
    }

    if p.IsIpV6 {
        fields["is_ipv6"] = "true"
    } else {
        fields["is_ipv6"] = "false"
    }

    return fields
}

func (p *Peer) Save() error {
    fields := p.GetFields()
    return InsertRow("peers", fields)
}

func (p *Peer) Update() error {
    fields := p.GetFields()
    return UpdateRow("peers", fields, "id=" + strconv.Itoa(p.Id))
}

func FindAvailablePeers(peerId string, infoHash string, isIpV6 bool) (peers []Peer, err error) {
  var ipV6 string
  if isIpV6 {
    ipV6 = "true"
  } else {
    ipV6 = "false"
  }

  sql := "SELECT * FROM peers WHERE peer_id!='" + Sanitize(peerId) + "' AND info_hash='" + Sanitize(infoHash) + "' AND is_ipv6=" + ipV6

  rows, err := Database.Query(sql)
  if err != nil {
    return peers, err
  }

  return GetPeersFromRows(rows)
}

func FindPeerByPeerIdAndInfoHash(peerId string, infoHash string) (p *Peer, err error) {
    sql := "SELECT * FROM peers WHERE peer_id='" + Sanitize(peerId) + "' AND info_hash='" + Sanitize(infoHash) + "'"

    rows, err := Database.Query(sql)
    if err != nil {
        return p, err
    }

    peers, err := GetPeersFromRows(rows)

    if err != nil || len(peers) == 0 {
        return p, err
    }

    return &peers[0], err
}

func GetPeersFromRows(rows *sql.Rows) (peers []Peer, err error) {
    for rows.Next() {
        var (
            id          int
            peerId      string
            ip          string
            port        int
            infoHash    string
            isIpV6      bool
        )

        err = rows.Scan(&id, &peerId, &ip, &port, &infoHash, &isIpV6)
        if err != nil {
            return peers, err
        }

        peer := new(Peer)
        peer.Id = id
        peer.PeerId = peerId
        peer.Ip = ip
        peer.Port = port
        peer.InfoHash = infoHash
        peer.IsIpV6 = isIpV6

        peers = append(peers, *peer)
    }

    return peers, err
}
