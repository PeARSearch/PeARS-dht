package peer

type PeerConfig struct {
	ID         int
	ListenPort int
	Global     bool
	Seed       int
	Target     string
}

func NewPeerConfig() PeerConfig {
	return PeerConfig{}
}

func (p *PeerConfig) Bootstrap() {
	// TODO this will Bootstrap the local peer
}
