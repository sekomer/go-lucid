package node

type NodeType string

const (
	FullNode NodeType = "full"
	BootNode NodeType = "boot"
	DevNode  NodeType = "dev"
)

type RpcConfig struct {
	Enabled bool     `yaml:"enabled"`
	Port    int      `yaml:"port"`
	Cors    string   `yaml:"cors"`
	Apis    []string `yaml:"apis"`
}

type P2pConfig struct {
	ListenPort  int  `yaml:"listen_port"`
	MinPeers    int  `yaml:"min_peers"`
	MaxPeers    int  `yaml:"max_peers"`
	GracePeriod int  `yaml:"grace_period"`
	Discovery   bool `yaml:"discovery"`
}

type DataConfig struct {
	Dir         string `yaml:"dir"`
	AutoMigrate bool   `yaml:"auto_migrate"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

type SyncConfig struct {
	Mode string `yaml:"mode"`
}

type MiningConfig struct {
	Enabled  bool   `yaml:"enabled"`
	MinerUrl string `yaml:"miner_url"`
}

type GenesisConfig struct {
	File string `yaml:"file"`
}

type DebugConfig struct {
	Enabled bool   `yaml:"enabled"`
	Seed    uint64 `yaml:"seed"`
	Peer    string `yaml:"peer"`
	Port    int    `yaml:"port"`
}

type NodeConfig struct {
	Id      string        `yaml:"id"`
	Type    NodeType      `yaml:"type"`
	Network string        `yaml:"network"`
	Rpc     RpcConfig     `yaml:"rpc"`
	P2p     P2pConfig     `yaml:"p2p"`
	Data    DataConfig    `yaml:"data"`
	Logging LoggingConfig `yaml:"logging"`
	Sync    SyncConfig    `yaml:"sync"`
	Mining  MiningConfig  `yaml:"mining"`
	Genesis GenesisConfig `yaml:"genesis"`
	Peers   []string      `yaml:"peers"`
	Debug   DebugConfig   `yaml:"debug"`
}

type FullNodeConfig struct {
	Node NodeConfig `yaml:"node"`
}
