package fn

type Protocol struct {
	Groupname string `json:"Groupname"`
	Protocol  string `json:"protocol"`
	RelayAddr string `json:"relay_addr"`
	As        string `json:"as"`
}

type IDFile struct {
	PeerId []string `json:"peerID"`
}

type IdPeer struct {
	PeerId string `json:"peerID"`
	Priv   string `json:"priv"`
}