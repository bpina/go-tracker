package data

import (
	"github.com/jackc/pgx"
)

type Torrent struct {
	InfoHash   []byte
	Incomplete int32
	Complete   int32
}

func (t *Torrent) Save() (bool, error) {
	pgx := "INSERT INTO torrents (info_hash, complete, incomplete) VALUES ($1, $2, $3)"

	commandTag, err := Database.Exec(pgx, []byte(t.InfoHash), t.Complete, t.Incomplete)
	if err != nil {
		return false, err
	}

	return commandTag.RowsAffected() > 0, err
}

func (t *Torrent) Update() (bool, error) {
	pgx := "UPDATE torrents set complete = $1, incomplete = $2 WHERE info_hash = $3"
	commandTag, err := Database.Exec(pgx, t.Complete, t.Incomplete, []byte(t.InfoHash))
	if err != nil {
		return false, err
	}

	return commandTag.RowsAffected() > 0, err
}

func (t *Torrent) Adjust(numWant int) {
	if numWant == 0 {
		t.Complete += 1
		if t.Incomplete > 0 {
			t.Incomplete -= 1
		}
	} else {
		t.Incomplete += 1
		if t.Complete > 0 {
			t.Complete -= 1
		}
	}
}

func (t *Torrent) GetPeers(peer *Peer) ([]Peer, error) {
	peers, err := FindAvailablePeers(peer.PeerId, t.InfoHash, peer.IsIpV6)
	return peers, err
}

func FindTorrent(infoHash []byte) (t *Torrent, err error) {
	pgx := "SELECT * FROM torrents WHERE info_hash = $1"

	rows, err := Database.Query(pgx, infoHash)
	if err != nil {
		return t, err
	}

	torrents, err := GetTorrentsFromRows(rows)
	if err != nil || len(torrents) == 0 {
		return t, err
	}

	rows.Close()

	return &torrents[0], err
}

func GetTorrentsFromRows(rows *pgx.Rows) (torrents []Torrent, err error) {
	for rows.Next() {
		var (
			info_hash  []byte
			complete   int32
			incomplete int32
		)

		err = rows.Scan(&info_hash, &incomplete, &complete)

		if err != nil {
			return torrents, err
		}

		torrent := new(Torrent)
		torrent.InfoHash = info_hash
		torrent.Complete = complete
		torrent.Incomplete = incomplete

		torrents = append(torrents, *torrent)
	}

	return torrents, err
}
