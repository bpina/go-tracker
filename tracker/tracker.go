package tracker

import (
    "net/url"
)

type Announce struct {
    InfoHash        string
    PeerId          string
    Port            string
    Uploaded        string
    Left            string
    Ip              string
    NumWant         string
    Event           string

    Data            url.Values
}

func ValidAnnounce(data url.Values) bool {
    valid := true
    fields := [...]string{"info_hash", "peer_id", "port", "uploaded",
                         "left", "numwant"}

    for i := range fields {
        if data.Get(fields[i]) == "" {
            valid = false
            break
        }
    }

    //TODO: implement stronger validation on available data. 

    return valid
}

func NewAnnounce(data url.Values) *Announce {
    if !ValidAnnounce(data) {
        return nil
    }

    announce := new(Announce)
    announce.InfoHash = data.Get("info_hash")
    announce.PeerId = data.Get("peer_id")
    announce.Port = data.Get("port")
    announce.Uploaded = data.Get("uploaded")
    announce.Left = data.Get("left")
    announce.Ip = data.Get("ip")
    announce.NumWant = data.Get("numwant")
    announce.Event = data.Get("event")
    announce.Data = data

    return announce
}


