package thp

import (
    "net/url"
    "strconv"
    "github.com/bpina/go-tracker/tools"
    "github.com/zeebo/bencode"
)

type Announce struct {
    InfoHash        string
    PeerId          string
    Port            int
    Uploaded        int
    Downloaded      int
    Left            int
    Ip              string
    NumWant         int
    Event           string

    Data            url.Values
}

func GetAnnounceErrors(data url.Values) []string {
    var errors []string
    fields := []string{"info_hash", "peer_id", "port", "uploaded",
                         "downloaded", "left"}

    for i := range fields {
        if data.Get(fields[i]) == "" {
           errors = append(errors, "Invalid value:" + fields[i])
        }
    }

    if len(errors) > 0 {
        return errors
    }

    if len(data.Get("info_hash")) != 20 {
        errors = append(errors, "Invalid value: info_hash")
    }

    if len(data.Get("peer_id")) != 20 {
        errors = append(errors, "Invalid value: peer_id")
    }

    i, err := strconv.Atoi(data.Get("port"))
    if err == nil {
        if i == 0 || i > 65535 {
            errors = append(errors, "Invalid value: port")
        }
    } else {
        errors = append(errors, "Invalid value: port")
    }

    i, err = strconv.Atoi(data.Get("uploaded"))
    if err != nil {
        errors = append(errors, "Invalid value: uploaded")
    }

    i, err = strconv.Atoi(data.Get("left"))
    if err != nil {
        errors = append(errors, "Invalid value: left")
    }

    i, err = strconv.Atoi(data.Get("downloaded"))
    if err != nil {
        errors = append(errors, "Invalid value: downloaded")
    }

    numwant := data.Get("numwant")
    if len(numwant) > 0 {
        i, err = strconv.Atoi(numwant)
        if err != nil {
            errors = append(errors, "Invalid value: numwant")
        }
    }

    eventMatch := false
    event := data.Get("event")
    if event != "" {
        for _, x := range &[...]string{"", "started", "stopped", "completed"} {
                if event == x {
                eventMatch = true
                break
            }
        }

        if !eventMatch {
            errors = append(errors, "Invalid value: event")
        }
    }

    //TODO: implement stronger validation on available data. 

    return errors
}


func NewAnnounce(data url.Values) (*Announce, []string) {
    errors := GetAnnounceErrors(data)

    if len(errors) > 0 {
        return nil, errors
    }

    announce := new(Announce)
    announce.InfoHash = data.Get("info_hash")
    announce.PeerId = data.Get("peer_id")
    announce.Port = tools.IntOrDefault(data.Get("port"))
    announce.Uploaded = tools.IntOrDefault(data.Get("uploaded"))
    announce.Left = tools.IntOrDefault(data.Get("left"))
    announce.Ip = data.Get("ip")
    announce.NumWant = tools.IntOrDefault(data.Get("numwant"))
    announce.Event = data.Get("event")
    announce.Data = data

    return announce, nil
}

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
