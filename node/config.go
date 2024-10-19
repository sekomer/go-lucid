package node

type NodeType string

const (
	FullNode NodeType = "full"
	BootNode NodeType = "boot"
	DevNode  NodeType = "dev"
)

type FullNodeConfig struct {
	Node struct {
		Id      string   `yaml:"id"`
		Type    NodeType `yaml:"type"`
		Network string   `yaml:"network"`
		Rpc     struct {
			Enabled bool     `yaml:"enabled"`
			Port    int      `yaml:"port"`
			Cors    string   `yaml:"cors"`
			Apis    []string `yaml:"apis"`
		} `yaml:"rpc"`
		P2p struct {
			ListenPort  int  `yaml:"listen_port"`
			MinPeers    int  `yaml:"min_peers"`
			MaxPeers    int  `yaml:"max_peers"`
			GracePeriod int  `yaml:"grace_period"`
			Discovery   bool `yaml:"discovery"`
		} `yaml:"p2p"`
		Data struct {
			Dir string `yaml:"dir"`
		} `yaml:"data"`
		Logging struct {
			Level string `yaml:"level"`
			File  string `yaml:"file"`
		} `yaml:"logging"`
		Sync struct {
			Mode string `yaml:"mode"`
		} `yaml:"sync"`
		Mining struct {
			Enabled  bool   `yaml:"enabled"`
			MinerUrl string `yaml:"miner_url"`
		} `yaml:"mining"`
		Genesis struct {
			File string `yaml:"file"`
		} `yaml:"genesis"`
		Peers []string `yaml:"peers"`
		Debug struct {
			Enabled bool   `yaml:"enabled"`
			Seed    uint64 `yaml:"seed"`
			Peer    string `yaml:"peer"`
			Port    int    `yaml:"port"`
		} `yaml:"debug"`
	} `yaml:"node"`
}
