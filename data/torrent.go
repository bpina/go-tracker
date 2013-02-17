package data

import (
    _ "github.com/bmizerany/pq"
    "database/sql"
    "log"
    "strconv"
)

type Torrent struct {
    InfoHash        string
    Incomplete      int
    Complete        int
}

func (t *Torrent) GetFields() map[string] string {
    fields := map[string] string {
        "info_hash": "'" + Sanitize(t.InfoHash) + "'",
        "complete": strconv.Itoa(t.Complete),
        "incomplete": strconv.Itoa(t.Incomplete),
    }

    return fields
}

func (t *Torrent) Save() error {
    //TODO: implement a one stop save/update
    fields := t.GetFields()
    return InsertRow("torrents", fields)
}

func (t *Torrent) Update() error {
    fields := t.GetFields()
    return UpdateRow("torrents", fields, "info_hash='" + Sanitize(t.InfoHash) + "'")
}

func FindTorrent(infoHash string) (t *Torrent, err error) {
    sanitized_info_hash := Sanitize(infoHash)
    sql := "SELECT * FROM torrents WHERE info_hash='" + sanitized_info_hash + "'"
    log.Printf(sql)

    rows, err := Database.Query(sql)
    if err != nil {
        return t, err
    }

    torrents, err := GetTorrentsFromRows(rows)
    if err != nil || len(torrents) == 0 {
        return t, err
    }

    return &torrents[0], err
}

func GetTorrentsFromRows(rows *sql.Rows) (torrents []Torrent, err error) {
    for rows.Next() {
        var (
            info_hash string
            incomplete int
            complete int
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
