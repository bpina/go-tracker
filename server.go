package main



import (
    "net/http"
    "log"
    "github.com/bpina/go-tracker/thp"
    "github.com/bpina/go-tracker/tools"
)

func AnnounceHandler(w http.ResponseWriter, req *http.Request) {
    log.Printf(req.URL.RawQuery)
    w.Header().Set("Content-Type", "text/plain")

    if req.Method == "POST" {
        response := thp.NewErrorResponse("Unsupported HTTP method.")
        w.Write([]byte(response.String()))
        return
    }

    req.ParseForm()
    announce, err := thp.NewAnnounce(req.Form)
    if err != nil {
        message := tools.FormatErrors(err)
        response := thp.NewErrorResponse(message)
        w.Write([]byte(response.String()))
        return
    } else {
        response := new(thp.Response)
        response.Interval = 30
        response.Complete = 999
        response.Incomplete = 999

        peer := new(thp.ConnectedPeer)
        peer.Ip = "192.168.3.99"
        peer.Port = 1337
        peer.PeerId = announce.PeerId

        response.Peers = append(response.Peers, *peer)

        message := response.String()
        log.Printf(message)
        w.Write([]byte(message))
       return
    }
}

func main() {
    http.HandleFunc("/announce", AnnounceHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
