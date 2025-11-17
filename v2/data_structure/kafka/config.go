package kafka

type Config struct {
	Name             string `toml:"name" json:"name"`
	Broker           string `toml:"broker" json:"broker"`
	Topic            string `toml:"topic" json:"topic"`
	Group            string `toml:"group" json:"group"`
	SecurityProtocol string `toml:"security.protocol" toml:"security.protocol"`
	SslCaLocation    string `toml:"ssl.ca.location" toml:"ssl.ca.location"`
	SaslMechanism    string `toml:"sasl.mechanism" toml:"sasl.mechanism"`
	SaslUsername     string `toml:"sasl.username" toml:"sasl.username"`
	SaslPassword     string `toml:"sasl.password" toml:"sasl.password"`
}
