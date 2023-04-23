package mongo

type Config struct {
	Address    string
	Username   string
	Password   string
	Database   string
	AuthSource string
	Timeout    int
}

type Configs map[string]*Config

func (m Configs) Use(key string) *Config {
	if item, ok := m[key]; ok {
		return item
	}
	return nil
}

type Reference struct {
	Key    string
	Config *Config
}

func (r *Reference) Parse(configs Configs) {
	r.Config = configs.Use(r.Key)
}
