package torrentfile

import (
	"bittorrent-client/peers"
	"github.com/jackpal/bencode-go"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

/*
	failure reason
	该关键字对应的值是一个可以读懂的字符串，指明GET请求失败的原因，如果返回信息中含有这个关键字，就不会再包含其他任何关键字

	warnging message
	该关键字对应的值是一个可以读懂的警告字符串

	interval
	指明客户端在下一次连接Tracker前所需等待的时间，以秒为单位

	min interval
	指明客户端在下一次连接Tracker前所需等待的最少时间，以秒为单位

	tracker id
	指明Tracker的ID

	complete
	一个整数，指明当前有多少个peer已经完成了整个共享文件的下载

	incomplete
	一个整数，指明当前有多少个peer还没有完成共享文件的下载

	peers
	返回各个peer的IP和端口号，它的值是一个字符串。首先是第一个peer的IP地址，然后是其端口号；
	接着是第二个peer的IP地址，然后是其端口号；依此类推
*/

// simple tracker server response
// example: d8:completei100e10:incompletei200e8:intervali1800e5:peers300:......e
type bencodeTrackerResp struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func (t *TorrentFile) buildTrackerURL(peerID [20]byte, port uint16) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}
	params := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(Port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Length)},
	}

	base.RawQuery = params.Encode()
	return base.String(), nil
}

func (t *TorrentFile) requestPeers(peerID [20]byte, port uint16) ([]peers.Peer, error) {
	url, err := t.buildTrackerURL(peerID, port)
	if err != nil {
		return nil, err
	}

	c := &http.Client{Timeout: 15 * time.Second}
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	trackerResp := bencodeTrackerResp{}
	err = bencode.Unmarshal(resp.Body, trackerResp)
	if err != nil {
		return nil, err
	}

	return peers.Unmarshal([]byte(trackerResp.Peers))
}
