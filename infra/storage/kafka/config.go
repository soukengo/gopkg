package kafka

type Config struct {
	Brokers     []string
	TopicPrefix string
}

type ConsumerConfig struct {
	Topics  []string
	Workers int
	GroupId string
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
