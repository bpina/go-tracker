package thp

import (
	"github.com/bpina/go-tracker/data"
	"github.com/bpina/go-tracker/tools"
	"log"
	"net"
	"net/http"
)

type Tracker struct {
	RemoteHost string
	IsIpV6     bool
	Request    *http.Request
}

func NewTracker(req *http.Request) (t *Tracker, err error) {
	err = ValidateTrackerRequest(req)
	if err != nil {
		return t, err
	}

	ip, err := ParseRemoteAddress(req.RemoteAddr)
	if err != nil {
		return t, err
	}

	log.Printf("parse IP: %s", string(ip))

	t = new(Tracker)
	t.RemoteHost = ip.String()
	t.IsIpV6 = ip.To4() == nil
	t.Request = req

	return t, err
}

func (t *Tracker) Execute() *Response {
	announce, announceErr := NewAnnounce(t.Request)
	if announceErr != nil {
		message := tools.FormatErrors(announceErr)
		return NewErrorResponse(message)
	}

	torrent, err := data.FindTorrent(announce.InfoHash)
	if err != nil {
		log.Print(err)
		return NewDatabaseErrorResponse()
	}

	if torrent == nil {
		torrent = &data.Torrent{announce.InfoHash, 0, 0}
		_, err := torrent.Save()
		if err != nil {
			log.Print(err)
			return NewErrorResponse("Could not find or create locate torrent.")
		}
	}

	torrent.Adjust(announce.NumWant)
	_, err = torrent.Update()
	if err != nil {
		log.Print(err)
		return NewDatabaseErrorResponse()
	}

	peer, err := data.FindPeerByPeerIdAndInfoHash(announce.PeerId, torrent.InfoHash)
	if err != nil {
		log.Print(err)
		return NewDatabaseErrorResponse()
	}

	if peer == nil {
		peer = t.CreatePeer(announce)
	} else {
		t.UpdatePeer(peer, announce)
	}

	peers, err := torrent.GetPeers(peer)
	if err != nil {
		log.Print(err)
		return NewDatabaseErrorResponse()
	}

	return NewTorrentResponse(torrent, peers)
}

func (t *Tracker) CreatePeer(announce *Announce) *data.Peer {
	log.Print("assigning IP to peer: %s", t.RemoteHost)
	peer := new(data.Peer)
	peer.PeerId = []byte(announce.PeerId)
	peer.Ip = t.RemoteHost
	peer.Port = int32(announce.Port)
	peer.InfoHash = []byte(announce.InfoHash)
	peer.IsIpV6 = t.IsIpV6
	peer.Save()

	return peer
}

func (t *Tracker) UpdatePeer(peer *data.Peer, announce *Announce) {
	peer.Port = int32(announce.Port)
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
