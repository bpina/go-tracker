package data

import (
	"github.com/jackc/pgx"
	"log"
)

type Peer struct {
	Id       int32
	PeerId   []byte
	Ip       string
	Port     int32
	InfoHash []byte
	IsIpV6   bool
}

func (p *Peer) Save() (bool, error) {
	pgx := "INSERT INTO peers (peer_id, ip, port, info_hash, is_ipv6) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	var id int32
	err := Database.QueryRow(pgx, []byte(p.PeerId), p.Ip, p.Port, p.InfoHash, p.IsIpV6).Scan(&id)

	if err != nil {
		log.Print(err)
		return false, err
	}

	p.Id = id

	return true, err
}

func (p *Peer) Update() (bool, error) {
	pgx := "UPDATE peers peer_id = $1, ip = $2, port = $3, info_hash = $4, is_ipv6 = $5)"

	commandTag, err := Database.Exec(pgx, []byte(p.PeerId), p.Ip, p.Port, []byte(p.InfoHash), p.IsIpV6)

	if err != nil {
		return false, err
	}

	return commandTag.RowsAffected() > 0, err

}

func FindAvailablePeers(peerId []byte, infoHash []byte, isIpV6 bool) (peers []Peer, err error) {
	pgx := "SELECT * FROM peers WHERE peer_id != $1 AND info_hash = $2 AND is_ipv6 = $3"

	log.Printf("\nfinding peers:\n peer_id: %x\n info_hash %x\n", peerId, infoHash)
	log.Printf("\nsql: %s\n", pgx)

	rows, err := Database.Query(pgx, peerId, infoHash, isIpV6)
	if err != nil {
		return peers, err
	}

	peers, err = GetPeersFromRows(rows)

	log.Printf("\npeer length: %i\n", len(peers))

	rows.Close()

	return peers, err
}

func FindPeerByPeerIdAndInfoHash(peerId []byte, infoHash []byte) (p *Peer, err error) {
	pgx := "SELECT * FROM peers WHERE peer_id = $1 AND info_hash = $2"

	rows, err := Database.Query(pgx, peerId, infoHash)
	if err != nil {
		log.Print(err)
		return p, err
	}

	peers, err := GetPeersFromRows(rows)

	if err != nil || len(peers) == 0 {
		return p, err
	}

	return &peers[0], err
}

func GetPeersFromRows(rows *pgx.Rows) (peers []Peer, err error) {
	for rows.Next() {
		var (
			id       int32
			peerId   []byte
			ip       string
			port     int32
			infoHash []byte
			isIpV6   bool
		)

		err = rows.Scan(&id, &peerId, &ip, &port, &infoHash, &isIpV6)
		if err != nil {
			log.Print(err)
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
