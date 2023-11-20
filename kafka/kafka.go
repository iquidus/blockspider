package kafka

type EventsConfig struct {
	Broker    string   `json:"broker"`
	Topic     string   `json:"topic"`
	Addresses []string `json:"addresses"`
	Topics    []string `json:"topics"`
}

type BlocksConfig struct {
	Broker string `json:"broker"`
	Topic  string `json:"topic"`
}

type Config struct {
	Blocks BlocksConfig   `json:"blocks"`
	Events []EventsConfig `json:"events"`
}
