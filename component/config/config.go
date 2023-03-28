package config

const (
	File  = "file"
	Nacos = "nacos"
)

type Config struct {
	Env    string
	Engine string
	File   FileConfig
	Nacos  NacosConfig
}
