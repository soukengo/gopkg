package discovery

type Config struct {
	Engine    string
	AppId     string           `yaml:"appId,omitempty"`
	Zookeeper *ZookeeperConfig `yaml:"zookeeper,omitempty"`
	Etcd      *EtcdConfig      `yaml:"etcd,omitempty"`
	Nacos     *NacosConfig     `yaml:"nacos,omitempty"`
}

type ZookeeperConfig struct {
	Nodes    []string
	RootPath string
}

type EtcdConfig struct {
	Nodes    []string
	Username string
	Password string
}

type NacosConfig struct {
	Server NacosServerConfig
	Client NacosClientConfig
}
type NacosServerConfig struct {
	Addr string
	Port uint64
}
type NacosClientConfig struct {
	GroupId   string
	Namespace string
	RegionId  string // the regionId for kms
	AccessKey string // the AccessKey for kms
	SecretKey string // the SecretKey for kms
	Username  string
	Password  string
	TimeoutMs uint64
	LogDir    string
	CacheDir  string
	LogLevel  string
}
