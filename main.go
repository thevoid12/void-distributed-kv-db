package main

import (
	"flag"
	"log"
	"main/config"
	"main/handler"
	"main/pkg/db"
	"net/http"

	"github.com/spf13/viper"
)

var (
	dbLocation = flag.String("db-location", "", "The path to the bolt db database")
	httpAddr   = flag.String("http-addr", "127.0.0.1:8080", "HTTP host and port")
	shard      = flag.String("shard", "", "The name of the shard for the data")
	configFile = flag.String("config-file", "config/config.json", "Config file for static sharding")
)

func parseFlags() {
	flag.Parse()
	if *dbLocation == "" {
		log.Fatalf("Must provide db-location")
	}
	if *shard == "" {
		log.Fatalf("Must provide shard name")
	}

}

func main() {
	parseFlags()

	log.Println("Reading configuration from:", *configFile)
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("config/") // path to look for the config file in

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var shardCount int
	var shardIdx int = -1
	var addrs = make(map[int]string)

	var config []config.Shard
	err = viper.UnmarshalKey("shards", &config)
	if err != nil {
		log.Fatalf("Error unmarshalling shard config: %v", err)
	}

	for _, s := range config {
		log.Printf("Available shard: %s, address: %s", s.Name, s.Address)
		addrs[s.Idx] = s.Address
		if s.Name == *shard {
			shardIdx = s.Idx
		}
	}

	shardCount = len(config)
	if shardIdx < 0 {
		log.Fatalf("Shard %q not found", *shard)
	}

	log.Printf("Shard count is %d, current shard: %d", shardCount, shardIdx)

	db, close, err := db.NewDatabase(*dbLocation)
	if err != nil {
		log.Fatalf("Error creating database: %v", err)
	}
	defer close()

	srv := handler.NewServer(db, shardIdx, shardCount, addrs)

	// Register handlers
	http.HandleFunc("/get", srv.GetHandler)
	http.HandleFunc("/set", srv.SetHandler)

	log.Printf("Starting server on %s", *httpAddr)
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
