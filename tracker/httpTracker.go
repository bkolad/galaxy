package tracker

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	i "github.com/bkolad/gTorrent/init"
	"github.com/bkolad/gTorrent/torrent"
)

type httpTracker struct {
	url string
}

//NewTracker creates default tracker
func NewTracker(info *torrent.Info, initState i.State, conf i.Configuration) (Tracker, error) {
	url, err := prepareURL(initState, conf.PeerID, conf.Port, info)
	if err != nil {
		return nil, err
	}
	return &httpTracker{url}, nil
}

func (t *httpTracker) Peers() ([]*torrent.PeerInfo, error) {
	resp, err := http.Get(t.url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("Can't reach http tracker " + t.url)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	bodyStr := string(body)
	dec := torrent.NewTrackerRspDecoder(bodyStr)
	rsp, err := dec.Decode()
	if err != nil {
		return nil, errors.New(bodyStr + " | " + err.Error())
	}
	return rsp.PeersInfo, nil
}

func prepareURL(initState i.State, peerID string, port int, info *torrent.Info) (string, error) {
	baseURL, err := url.Parse(info.Announce)
	if err != nil {
		return "", errors.New("Malformed URL: " + err.Error())
	}

	params := url.Values{}
	params.Add("info_hash", "lol") //string(info.InfoHash))
	params.Add("peer_id", peerID)
	params.Add("port", strconv.Itoa(port))
	params.Add("compact", "1")
	params.Add("event", "started")
	params.Add("uploaded", strconv.Itoa(initState.Uploaded))
	params.Add("downloaded", strconv.Itoa(initState.Downloaded))
	params.Add("left", strconv.Itoa(initState.Left))
	baseURL.RawQuery = params.Encode()
	return baseURL.String(), nil
}
