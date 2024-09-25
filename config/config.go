package config

// Shard describes a shard that holds the appropriate set of keys.
// Each shard has unique set of keys.
type Shard struct {
	Name    string `mapstructure:"name"`
	Idx     int    `mapstructure:"idx"`
	Address string `mapstructure:"address"`
}

// Config describes the sharding config.
type Config struct {
	Shards []Shard
}
