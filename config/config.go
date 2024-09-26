package config

import (
	"fmt"
	"hash/fnv"
)

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

// Shards represents an easier-to-use representation of
// the sharding config: the shards count, current index and
// the addresses of all other shards too.
type Shards struct {
	Count     int
	CurIdx    int
	Addrs     map[int]string
	ShardName map[int]string
}

// ParseShards converts and verifies the list of shards
// specified in the config into a form that can be used
// for routing.
func ParseShards(shards []Shard, curShardName string) (*Shards, error) {
	shardCount := len(shards)
	shardIdx := -1
	addrs := make(map[int]string)
	name := make(map[int]string)
	for _, s := range shards {
		if _, ok := addrs[s.Idx]; ok {
			return nil, fmt.Errorf("duplicate shard index: %d", s.Idx)
		}
		addrs[s.Idx] = s.Address
		name[s.Idx] = s.Name
		if s.Name == curShardName {
			shardIdx = s.Idx
		}
	}
	for i := 0; i < shardCount; i++ {
		if _, ok := addrs[i]; !ok {
			return nil, fmt.Errorf("shard %d is not found", i)
		}
	}
	if shardIdx < 0 {
		return nil, fmt.Errorf("shard %q was not found", curShardName)
	}
	return &Shards{
		Count:     shardCount,
		CurIdx:    shardIdx,
		Addrs:     addrs,
		ShardName: name,
	}, nil
}

// Index returns the shard number for the corresponding key.
func (s *Shards) Index(key string) int {
	h := fnv.New64()
	h.Write([]byte(key))
	return int(h.Sum64() % uint64(s.Count))
}
